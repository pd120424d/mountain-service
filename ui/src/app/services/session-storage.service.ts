import { Injectable } from '@angular/core';
import { LoggingService } from './logging.service';

export interface SessionImageData {
  preview: string;
  fileName: string;
  fileSize: number;
  fileType: string;
  uploadResult?: {
    blobUrl: string;
    blobName: string;
    size: number;
    message: string;
  };
  timestamp: number;
}

@Injectable({
  providedIn: 'root'
})
export class SessionStorageService {
  private readonly IMAGE_CACHE_PREFIX = 'mountain_service_image_';
  private readonly CACHE_DURATION = 24 * 60 * 60 * 1000; // 24 hours

  constructor(private logger: LoggingService) {
    this.logger.info('SessionStorageService initialized');
    this.cleanupExpiredImages();
  }

  /**
   * Store image data in session storage
   */
  storeImageData(key: string, imageData: SessionImageData): void {
    try {
      const storageKey = this.getStorageKey(key);
      const dataToStore = {
        ...imageData,
        timestamp: Date.now()
      };
      
      sessionStorage.setItem(storageKey, JSON.stringify(dataToStore));
      this.logger.debug(`Image data stored in session storage: ${key}`);
    } catch (error) {
      this.logger.error('Failed to store image data in session storage:', error);
      // If storage is full, try to clean up and retry
      this.cleanupExpiredImages();
      try {
        const storageKey = this.getStorageKey(key);
        sessionStorage.setItem(storageKey, JSON.stringify(imageData));
      } catch (retryError) {
        this.logger.error('Failed to store image data after cleanup:', retryError);
      }
    }
  }

  /**
   * Retrieve image data from session storage
   */
  getImageData(key: string): SessionImageData | null {
    try {
      const storageKey = this.getStorageKey(key);
      const storedData = sessionStorage.getItem(storageKey);
      
      if (!storedData) {
        return null;
      }

      const imageData: SessionImageData = JSON.parse(storedData);
      
      // Check if data is still valid
      if (this.isDataExpired(imageData)) {
        this.removeImageData(key);
        return null;
      }

      return imageData;
    } catch (error) {
      this.logger.error('Failed to retrieve image data from session storage:', error);
      return null;
    }
  }

  /**
   * Remove image data from session storage
   */
  removeImageData(key: string): void {
    try {
      const storageKey = this.getStorageKey(key);
      sessionStorage.removeItem(storageKey);
      this.logger.debug(`Image data removed from session storage: ${key}`);
    } catch (error) {
      this.logger.error('Failed to remove image data from session storage:', error);
    }
  }

  /**
   * Clear all image data from session storage
   */
  clearAllImageData(): void {
    try {
      const keysToRemove: string[] = [];
      
      for (let i = 0; i < sessionStorage.length; i++) {
        const key = sessionStorage.key(i);
        if (key && key.startsWith(this.IMAGE_CACHE_PREFIX)) {
          keysToRemove.push(key);
        }
      }

      keysToRemove.forEach(key => sessionStorage.removeItem(key));
      this.logger.info(`Cleared ${keysToRemove.length} image data entries from session storage`);
    } catch (error) {
      this.logger.error('Failed to clear image data from session storage:', error);
    }
  }

  /**
   * Get all stored image keys
   */
  getAllImageKeys(): string[] {
    const keys: string[] = [];
    
    try {
      for (let i = 0; i < sessionStorage.length; i++) {
        const key = sessionStorage.key(i);
        if (key && key.startsWith(this.IMAGE_CACHE_PREFIX)) {
          keys.push(key.replace(this.IMAGE_CACHE_PREFIX, ''));
        }
      }
    } catch (error) {
      this.logger.error('Failed to get image keys from session storage:', error);
    }

    return keys;
  }

  /**
   * Get storage usage information
   */
  getStorageInfo(): { totalKeys: number; estimatedSize: string; availableSpace: string } {
    let totalKeys = 0;
    let estimatedSize = 0;

    try {
      for (let i = 0; i < sessionStorage.length; i++) {
        const key = sessionStorage.key(i);
        if (key && key.startsWith(this.IMAGE_CACHE_PREFIX)) {
          totalKeys++;
          const value = sessionStorage.getItem(key);
          if (value) {
            estimatedSize += value.length * 2; // Rough estimate (UTF-16)
          }
        }
      }

      // Estimate available space (sessionStorage limit is usually 5-10MB)
      const estimatedLimit = 5 * 1024 * 1024; // 5MB
      const availableSpace = Math.max(0, estimatedLimit - estimatedSize);

      return {
        totalKeys,
        estimatedSize: this.formatBytes(estimatedSize),
        availableSpace: this.formatBytes(availableSpace)
      };
    } catch (error) {
      this.logger.error('Failed to get storage info:', error);
      return {
        totalKeys: 0,
        estimatedSize: '0 B',
        availableSpace: 'Unknown'
      };
    }
  }

  /**
   * Store user preferences for image handling
   */
  storeImagePreferences(preferences: {
    autoUpload?: boolean;
    compressionQuality?: number;
    maxFileSize?: number;
  }): void {
    try {
      sessionStorage.setItem('mountain_service_image_prefs', JSON.stringify(preferences));
      this.logger.debug('Image preferences stored');
    } catch (error) {
      this.logger.error('Failed to store image preferences:', error);
    }
  }

  /**
   * Get user preferences for image handling
   */
  getImagePreferences(): {
    autoUpload: boolean;
    compressionQuality: number;
    maxFileSize: number;
  } {
    try {
      const stored = sessionStorage.getItem('mountain_service_image_prefs');
      if (stored) {
        const preferences = JSON.parse(stored);
        return {
          autoUpload: preferences.autoUpload ?? true,
          compressionQuality: preferences.compressionQuality ?? 0.8,
          maxFileSize: preferences.maxFileSize ?? 5 * 1024 * 1024 // 5MB
        };
      }
    } catch (error) {
      this.logger.error('Failed to get image preferences:', error);
    }

    // Return defaults
    return {
      autoUpload: true,
      compressionQuality: 0.8,
      maxFileSize: 5 * 1024 * 1024 // 5MB
    };
  }

  private getStorageKey(key: string): string {
    return `${this.IMAGE_CACHE_PREFIX}${key}`;
  }

  private isDataExpired(imageData: SessionImageData): boolean {
    return (Date.now() - imageData.timestamp) > this.CACHE_DURATION;
  }

  private cleanupExpiredImages(): void {
    try {
      const keysToRemove: string[] = [];
      
      for (let i = 0; i < sessionStorage.length; i++) {
        const key = sessionStorage.key(i);
        if (key && key.startsWith(this.IMAGE_CACHE_PREFIX)) {
          const value = sessionStorage.getItem(key);
          if (value) {
            try {
              const imageData: SessionImageData = JSON.parse(value);
              if (this.isDataExpired(imageData)) {
                keysToRemove.push(key);
              }
            } catch (parseError) {
              // Invalid data, remove it
              keysToRemove.push(key);
            }
          }
        }
      }

      keysToRemove.forEach(key => sessionStorage.removeItem(key));
      
      if (keysToRemove.length > 0) {
        this.logger.info(`Cleaned up ${keysToRemove.length} expired image entries from session storage`);
      }
    } catch (error) {
      this.logger.error('Failed to cleanup expired images:', error);
    }
  }

  private formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }
}
