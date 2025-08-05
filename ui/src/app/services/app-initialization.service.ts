import { Injectable } from '@angular/core';
import { ImageUploadService } from './image-upload.service';
import { LoggingService } from './logging.service';

@Injectable({
  providedIn: 'root'
})
export class AppInitializationService {

  constructor(
    private imageUploadService: ImageUploadService,
    private logger: LoggingService
  ) {}

  /**
   * Initialize the application
   * This method should be called during app startup
   */
  initialize(): Promise<void> {
    return new Promise((resolve) => {
      this.logger.info('Initializing application...');

      try {
        // Preload cached images from session storage
        this.imageUploadService.preloadCachedImages();

        // Log storage info for debugging
        const storageInfo = this.imageUploadService.getSessionStorageInfo();
        this.logger.info(`Session storage info: ${storageInfo.totalKeys} images, ${storageInfo.estimatedSize} used, ${storageInfo.availableSpace} available`);

        // Log user preferences
        const preferences = this.imageUploadService.getImagePreferences();
        this.logger.debug('Image preferences:', preferences);

        this.logger.info('Application initialization completed');
        resolve();
      } catch (error) {
        this.logger.error('Error during application initialization:', error);
        // Don't fail the app startup, just log the error
        resolve();
      }
    });
  }

  /**
   * Cleanup resources before app shutdown
   */
  cleanup(): void {
    this.logger.info('Cleaning up application resources...');
    
    try {
      // Optional: Clear expired cache entries
      // Note: We don't clear all cache as user might want to keep images between sessions
      this.logger.info('Application cleanup completed');
    } catch (error) {
      this.logger.error('Error during application cleanup:', error);
    }
  }
}
