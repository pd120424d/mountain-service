package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service struct {
	Name       string `json:"name"`
	BaseURL    string `json:"baseURL"`
	HealthPath string `json:"healthPath"`
	SpecPath   string `json:"specPath"`
}

type Config struct {
	ExternalScheme   string
	ExternalHost     string
	ExternalBasePath string
	Services         []Service
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func loadConfig() (*Config, error) {
	extScheme := getenv("EXTERNAL_SCHEME", "https")
	extHost := getenv("EXTERNAL_HOST", "mountain-service.duckdns.org")
	extBase := getenv("EXTERNAL_BASE_PATH", "/api/v1")

	var services []Service
	if raw := os.Getenv("SERVICES_JSON"); strings.TrimSpace(raw) != "" {
		if err := json.Unmarshal([]byte(raw), &services); err != nil {
			return nil, fmt.Errorf("failed to parse SERVICES_JSON: %w", err)
		}
	} else {
		services = []Service{
			{Name: "employee", BaseURL: "http://employee-service:8082", HealthPath: "/api/v1/health", SpecPath: "/swagger.json"},
			{Name: "urgency", BaseURL: "http://urgency-service:8083", HealthPath: "/api/v1/health", SpecPath: "/swagger.json"},
			{Name: "activity", BaseURL: "http://activity-service:8084", HealthPath: "/api/v1/health", SpecPath: "/swagger.json"},
		}
	}

	return &Config{
		ExternalScheme:   extScheme,
		ExternalHost:     extHost,
		ExternalBasePath: extBase,
		Services:         services,
	}, nil
}

func checkHealth(basedURL, healthPath string) bool {
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	req, _ := http.NewRequest("GET", strings.TrimRight(basedURL, "/")+healthPath, nil)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func rewriteSpecBytes(spec []byte, scheme, host, basePath string) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(spec, &m); err != nil {
		return nil, fmt.Errorf("invalid JSON spec: %w", err)
	}
	if v, ok := m["swagger"]; ok {
		// Swagger 2.0
		_ = v
		m["host"] = host
		m["schemes"] = []string{scheme}
		if basePath != "" {
			m["basePath"] = basePath
		}
		return json.MarshalIndent(m, "", "  ")
	}
	if _, ok := m["openapi"]; ok {
		// OpenAPI 3.x
		url := fmt.Sprintf("%s://%s%s", scheme, host, basePath)
		m["servers"] = []map[string]string{{"url": url}}
		return json.MarshalIndent(m, "", "  ")
	}
	// Unknown, return as-is
	return spec, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Errorf("docs-aggregator config error: %w", err))
	}

	log, err := utils.NewLogger("docs-aggregator")
	if err != nil {
		panic("failed to create logger")
	}

	ctx, _ := utils.EnsureRequestID(context.Background())
	log.WithContext(ctx).Info("Starting Docs Aggregator service")
	defer utils.TimeOperation(log, "DocsAggregatorService.main")()

	r := gin.Default()
	r.Use(log.RequestLogger())

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})

	docs := r.Group("/docs")
	apiDocs := r.Group("/api/v1/docs")

	rootHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Docs aggregator is running",
			"endpoints": []string{
				"/docs/swagger-config.json",
				"/docs/specs/:service",
				"/api/v1/docs/swagger-config.json",
				"/api/v1/docs/specs/:service",
			},
		})
	}

	{
		docs.GET("", rootHandler)
		apiDocs.GET("", rootHandler)

		swaggerConfigHandler := func(c *gin.Context) {
			c.Header("Cache-Control", "public, max-age=60")
			reqLog := log.WithContext(c.Request.Context())
			urls := make([]map[string]string, 0, len(cfg.Services))

			basePath := "/docs/specs"
			if strings.HasPrefix(c.Request.URL.Path, "/api/v1/docs") {
				basePath = "/api/v1/docs/specs"
			}

			for _, s := range cfg.Services {
				if checkHealth(s.BaseURL, s.HealthPath) {
					reqLog.Debugf("service %s healthy for docs", s.Name)
					urls = append(urls, map[string]string{
						"name": cases.Title(language.English, cases.NoLower).String(s.Name) + " API",
						"url":  fmt.Sprintf("%s/%s.json", basePath, s.Name),
					})
				}
			}
			resp := map[string]interface{}{
				"urls":        urls,
				"deepLinking": true,
				"layout":      "BaseLayout",
			}
			c.JSON(http.StatusOK, resp)
		}

		docs.GET("/swagger-config.json", swaggerConfigHandler)
		apiDocs.GET("/swagger-config.json", swaggerConfigHandler)

		specsHandler := func(c *gin.Context) {
			reqLog := log.WithContext(c.Request.Context())
			name := c.Param("service")

			// Make sure to trim any blank spaces and .json extension if present
			name = strings.TrimSpace(strings.ToLower(name))
			name = strings.TrimSuffix(name, ".json")

			// Optional aliases
			switch name {
			case "employees":
				name = "employee"
			case "urgencies":
				name = "urgency"
			case "activities":
				name = "activity"
			}

			var svc *Service
			for i := range cfg.Services {
				if strings.EqualFold(cfg.Services[i].Name, name) {
					svc = &cfg.Services[i]
					break
				}
			}

			// Fallback, try to match by BaseURL (e.g., employee-service)
			if svc == nil {
				for i := range cfg.Services {
					if strings.Contains(strings.ToLower(cfg.Services[i].BaseURL), name) {
						svc = &cfg.Services[i]
						break
					}
				}
			}
			if svc == nil {
				reqLog.Errorf("unknown service for docs: %s", name)
				c.JSON(http.StatusNotFound, gin.H{"error": "unknown service"})
				return
			}

			client := &http.Client{Timeout: 5 * time.Second}
			upstream := strings.TrimRight(svc.BaseURL, "/") + svc.SpecPath
			req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, upstream, nil)
			resp, err := client.Do(req)
			if err != nil {
				reqLog.Errorf("failed to fetch upstream spec: %v", err)
				c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch upstream spec"})
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				reqLog.Errorf("upstream status %d for %s", resp.StatusCode, upstream)
				c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("upstream status %d", resp.StatusCode)})
				return
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				reqLog.Errorf("failed to read upstream spec: %v", err)
				c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read upstream spec"})
				return
			}
			rewritten, err := rewriteSpecBytes(b, cfg.ExternalScheme, cfg.ExternalHost, cfg.ExternalBasePath)
			if err != nil {
				// Fallback to original if rewrite fails
				rewritten = b
			}
			c.Header("Cache-Control", "public, max-age=60")
			c.Data(http.StatusOK, "application/json", rewritten)
		}

		docs.GET("/specs/:service", specsHandler)
		apiDocs.GET("/specs/:service", specsHandler)

	}

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	log.Infof("docs-aggregator listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
