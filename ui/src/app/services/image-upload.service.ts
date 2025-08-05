import { Injectable } from '@angular/core';
import { HttpClient, HttpEvent, HttpEventType, HttpProgressEvent, HttpResponse } from '@angular/common/http';
import { Observable, BehaviorSubject, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { LoggingService } from './logging.service';
import { SessionStorageService, SessionImageData } from './session-storage.service';

export interface ImageUploadProgress {
  progress: number;
  status: 'uploading' | 'completed' | 'error';
  message?: string;
}

export interface ImageUploadResult {
  blobUrl: string;
  blobName: string;
  size: number;
  message: string;
}

export interface CachedImage {
  file: File;
  preview: string;
  uploadResult?: ImageUploadResult;
  timestamp: number;
}

@Injectable({
  providedIn: 'root'
})
export class ImageUploadService {
  private baseApiUrl = environment.useMockApi
    ? '/api/v1' // Mock server URL
    : `${environment.apiUrl}`; // Real API

  private imageCache = new Map<string, CachedImage>();
  private readonly CACHE_DURATION = 24 * 60 * 60 * 1000; // 24 hours in milliseconds

  constructor(
    private http: HttpClient,
    private logger: LoggingService,
    private sessionStorage: SessionStorageService
  ) {
    this.logger.info('ImageUploadService initialized');
    this.cleanupExpiredCache();
  }

  /**
   * Upload profile picture for an employee
   */
  uploadProfilePicture(employeeId: number, file: File): Observable<ImageUploadProgress> {
    const formData = new FormData();
    formData.append('file', file);

    const uploadUrl = `${this.baseApiUrl}/employees/${employeeId}/profile-picture`;
    
    return this.http.post<ImageUploadResult>(uploadUrl, formData, {
      reportProgress: true,
      observe: 'events'
    }).pipe(
      map((event: HttpEvent<ImageUploadResult>) => this.mapHttpEventToProgress(event)),
      tap((progress) => {
        if (progress.status === 'completed') {
          this.logger.info(`Profile picture uploaded successfully for employee ${employeeId}`);
        }
      }),
      catchError((error) => {
        this.logger.error(`Failed to upload profile picture for employee ${employeeId}:`, error);
        return throwError(() => ({
          progress: 0,
          status: 'error' as const,
          message: this.getErrorMessage(error)
        }));
      })
    );
  }

  /**
   * Delete profile picture for an employee
   */
  deleteProfilePicture(employeeId: number, blobName: string): Observable<any> {
    const deleteUrl = `${this.baseApiUrl}/employees/${employeeId}/profile-picture?blobName=${encodeURIComponent(blobName)}`;
    
    return this.http.delete(deleteUrl).pipe(
      tap(() => {
        this.logger.info(`Profile picture deleted successfully for employee ${employeeId}`);
        this.removeFromCache(`employee-${employeeId}`);
      }),
      catchError((error) => {
        this.logger.error(`Failed to delete profile picture for employee ${employeeId}:`, error);
        return throwError(() => error);
      })
    );
  }

  /**
   * Validate image file before upload
   */
  validateImageFile(file: File): { valid: boolean; error?: string } {
    // Check file size (max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (file.size > maxSize) {
      return { valid: false, error: 'File size exceeds 5MB limit' };
    }

    // Check file type
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      return { valid: false, error: 'Invalid file type. Please select a JPEG, PNG, GIF, or WebP image.' };
    }

    return { valid: true };
  }

  /**
   * Create image preview URL
   */
  createImagePreview(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => resolve(e.target?.result as string);
      reader.onerror = () => reject(new Error('Failed to read file'));
      reader.readAsDataURL(file);
    });
  }

  /**
   * Cache image for session management
   */
  cacheImage(key: string, file: File, preview: string, uploadResult?: ImageUploadResult): void {
    // Store in memory cache
    this.imageCache.set(key, {
      file,
      preview,
      uploadResult,
      timestamp: Date.now()
    });

    // Store in session storage for persistence
    const sessionData: SessionImageData = {
      preview,
      fileName: file.name,
      fileSize: file.size,
      fileType: file.type,
      uploadResult,
      timestamp: Date.now()
    };
    this.sessionStorage.storeImageData(key, sessionData);

    this.logger.debug(`Image cached with key: ${key}`);
  }

  /**
   * Get cached image
   */
  getCachedImage(key: string): CachedImage | null {
    // First check memory cache
    const cached = this.imageCache.get(key);
    if (cached && this.isCacheValid(cached)) {
      return cached;
    }
    if (cached) {
      this.imageCache.delete(key); // Remove expired cache
    }

    // If not in memory, check session storage
    const sessionData = this.sessionStorage.getImageData(key);
    if (sessionData) {
      // Reconstruct cached image from session data
      // Note: We can't reconstruct the File object from session storage
      // So we'll return a partial cache that can be used for display
      return {
        file: new File([], sessionData.fileName, { type: sessionData.fileType }),
        preview: sessionData.preview,
        uploadResult: sessionData.uploadResult,
        timestamp: sessionData.timestamp
      };
    }

    return null;
  }

  /**
   * Remove image from cache
   */
  removeFromCache(key: string): void {
    this.imageCache.delete(key);
    this.sessionStorage.removeImageData(key);
    this.logger.debug(`Image removed from cache: ${key}`);
  }

  /**
   * Clear all cached images
   */
  clearCache(): void {
    this.imageCache.clear();
    this.sessionStorage.clearAllImageData();
    this.logger.info('Image cache cleared');
  }

  /**
   * Get all cached images (for debugging)
   */
  getCacheInfo(): { key: string; size: number; timestamp: Date }[] {
    return Array.from(this.imageCache.entries()).map(([key, value]) => ({
      key,
      size: value.file.size,
      timestamp: new Date(value.timestamp)
    }));
  }

  /**
   * Get session storage info
   */
  getSessionStorageInfo() {
    return this.sessionStorage.getStorageInfo();
  }

  /**
   * Preload cached images on app startup
   */
  preloadCachedImages(): void {
    const keys = this.sessionStorage.getAllImageKeys();
    let loadedCount = 0;

    keys.forEach(key => {
      const sessionData = this.sessionStorage.getImageData(key);
      if (sessionData) {
        // Create a minimal cached image entry for display purposes
        this.imageCache.set(key, {
          file: new File([], sessionData.fileName, { type: sessionData.fileType }),
          preview: sessionData.preview,
          uploadResult: sessionData.uploadResult,
          timestamp: sessionData.timestamp
        });
        loadedCount++;
      }
    });

    if (loadedCount > 0) {
      this.logger.info(`Preloaded ${loadedCount} cached images from session storage`);
    }
  }

  /**
   * Get user preferences for image handling
   */
  getImagePreferences() {
    return this.sessionStorage.getImagePreferences();
  }

  /**
   * Store user preferences for image handling
   */
  storeImagePreferences(preferences: {
    autoUpload?: boolean;
    compressionQuality?: number;
    maxFileSize?: number;
  }): void {
    this.sessionStorage.storeImagePreferences(preferences);
  }

  private mapHttpEventToProgress(event: HttpEvent<ImageUploadResult>): ImageUploadProgress {
    switch (event.type) {
      case HttpEventType.UploadProgress:
        const progressEvent = event as HttpProgressEvent;
        const progress = progressEvent.total 
          ? Math.round(100 * progressEvent.loaded / progressEvent.total)
          : 0;
        return {
          progress,
          status: 'uploading',
          message: `Uploading... ${progress}%`
        };

      case HttpEventType.Response:
        const response = event as HttpResponse<ImageUploadResult>;
        return {
          progress: 100,
          status: 'completed',
          message: response.body?.message || 'Upload completed successfully'
        };

      default:
        return {
          progress: 0,
          status: 'uploading',
          message: 'Preparing upload...'
        };
    }
  }

  private getErrorMessage(error: any): string {
    if (error?.error?.error) {
      return error.error.error;
    }
    if (error?.message) {
      return error.message;
    }
    if (error?.status === 0) {
      return 'Network error. Please check your connection.';
    }
    if (error?.status >= 500) {
      return 'Server error. Please try again later.';
    }
    return 'Upload failed. Please try again.';
  }

  private isCacheValid(cached: CachedImage): boolean {
    return (Date.now() - cached.timestamp) < this.CACHE_DURATION;
  }

  private cleanupExpiredCache(): void {
    const now = Date.now();
    for (const [key, value] of this.imageCache.entries()) {
      if ((now - value.timestamp) >= this.CACHE_DURATION) {
        this.imageCache.delete(key);
      }
    }
  }
}
