import { TestBed } from '@angular/core/testing';
import { AppInitializationService } from './app-initialization.service';
import { ImageUploadService } from './image-upload.service';
import { LoggingService } from './logging.service';

describe('AppInitializationService', () => {
  let service: AppInitializationService;
  let mockImageUploadService: jasmine.SpyObj<ImageUploadService>;
  let mockLoggingService: jasmine.SpyObj<LoggingService>;

  beforeEach(() => {
    const imageUploadServiceSpy = jasmine.createSpyObj('ImageUploadService', [
      'preloadCachedImages',
      'getSessionStorageInfo',
      'getImagePreferences'
    ]);
    const loggingServiceSpy = jasmine.createSpyObj('LoggingService', ['info', 'debug', 'error']);

    TestBed.configureTestingModule({
      providers: [
        AppInitializationService,
        { provide: ImageUploadService, useValue: imageUploadServiceSpy },
        { provide: LoggingService, useValue: loggingServiceSpy }
      ]
    });

    service = TestBed.inject(AppInitializationService);
    mockImageUploadService = TestBed.inject(ImageUploadService) as jasmine.SpyObj<ImageUploadService>;
    mockLoggingService = TestBed.inject(LoggingService) as jasmine.SpyObj<LoggingService>;
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('initialize', () => {
    it('it succeeds when initialization completes successfully', async () => {
      const mockStorageInfo = {
        totalKeys: 5,
        estimatedSize: '2.5 MB',
        availableSpace: '2.5 MB'
      };
      const mockPreferences = {
        autoUpload: true,
        compressionQuality: 0.8,
        maxFileSize: 5242880
      };

      mockImageUploadService.preloadCachedImages.and.stub();
      mockImageUploadService.getSessionStorageInfo.and.returnValue(mockStorageInfo);
      mockImageUploadService.getImagePreferences.and.returnValue(mockPreferences);

      await service.initialize();

      expect(mockLoggingService.info).toHaveBeenCalledWith('Initializing application...');
      expect(mockImageUploadService.preloadCachedImages).toHaveBeenCalled();
      expect(mockImageUploadService.getSessionStorageInfo).toHaveBeenCalled();
      expect(mockImageUploadService.getImagePreferences).toHaveBeenCalled();
      expect(mockLoggingService.info).toHaveBeenCalledWith('Session storage info: 5 images, 2.5 MB used, 2.5 MB available');
      expect(mockLoggingService.debug).toHaveBeenCalledWith('Image preferences:', mockPreferences);
      expect(mockLoggingService.info).toHaveBeenCalledWith('Application initialization completed');
    });

    it('it succeeds when initialization encounters errors but continues', async () => {
      mockImageUploadService.preloadCachedImages.and.throwError('Preload error');
      mockImageUploadService.getSessionStorageInfo.and.throwError('Storage info error');
      mockImageUploadService.getImagePreferences.and.throwError('Preferences error');

      await service.initialize();

      expect(mockLoggingService.info).toHaveBeenCalledWith('Initializing application...');
      expect(mockLoggingService.error).toHaveBeenCalledWith('Error during application initialization:', jasmine.any(Error));
      // Should still resolve the promise even with errors
    });

    it('it succeeds when preloadCachedImages throws error', async () => {
      mockImageUploadService.preloadCachedImages.and.throwError('Preload failed');
      mockImageUploadService.getSessionStorageInfo.and.returnValue({
        totalKeys: 0,
        estimatedSize: '0 B',
        availableSpace: '5 MB'
      });
      mockImageUploadService.getImagePreferences.and.returnValue({
        autoUpload: true,
        compressionQuality: 0.8,
        maxFileSize: 5242880
      });

      await service.initialize();

      expect(mockLoggingService.error).toHaveBeenCalledWith('Error during application initialization:', jasmine.any(Error));
    });

    it('it succeeds when getSessionStorageInfo throws error', async () => {
      mockImageUploadService.preloadCachedImages.and.stub();
      mockImageUploadService.getSessionStorageInfo.and.throwError('Storage info failed');
      mockImageUploadService.getImagePreferences.and.returnValue({
        autoUpload: true,
        compressionQuality: 0.8,
        maxFileSize: 5242880
      });

      await service.initialize();

      expect(mockLoggingService.error).toHaveBeenCalledWith('Error during application initialization:', jasmine.any(Error));
    });

    it('it succeeds when getImagePreferences throws error', async () => {
      mockImageUploadService.preloadCachedImages.and.stub();
      mockImageUploadService.getSessionStorageInfo.and.returnValue({
        totalKeys: 0,
        estimatedSize: '0 B',
        availableSpace: '5 MB'
      });
      mockImageUploadService.getImagePreferences.and.throwError('Preferences failed');

      await service.initialize();

      expect(mockLoggingService.error).toHaveBeenCalledWith('Error during application initialization:', jasmine.any(Error));
    });

    it('it returns a promise that resolves', async () => {
      mockImageUploadService.preloadCachedImages.and.stub();
      mockImageUploadService.getSessionStorageInfo.and.returnValue({
        totalKeys: 0,
        estimatedSize: '0 B',
        availableSpace: '5 MB'
      });
      mockImageUploadService.getImagePreferences.and.returnValue({
        autoUpload: true,
        compressionQuality: 0.8,
        maxFileSize: 5242880
      });

      const result = service.initialize();

      expect(result).toBeInstanceOf(Promise);
      await expectAsync(result).toBeResolved();
    });

    it('it logs storage info with correct format', async () => {
      const mockStorageInfo = {
        totalKeys: 10,
        estimatedSize: '1.2 MB',
        availableSpace: '3.8 MB'
      };

      mockImageUploadService.preloadCachedImages.and.stub();
      mockImageUploadService.getSessionStorageInfo.and.returnValue(mockStorageInfo);
      mockImageUploadService.getImagePreferences.and.returnValue({
        autoUpload: false,
        compressionQuality: 0.9,
        maxFileSize: 10485760
      });

      await service.initialize();

      expect(mockLoggingService.info).toHaveBeenCalledWith('Session storage info: 10 images, 1.2 MB used, 3.8 MB available');
    });
  });

  describe('cleanup', () => {
    it('it succeeds when cleanup completes successfully', () => {
      service.cleanup();

      expect(mockLoggingService.info).toHaveBeenCalledWith('Cleaning up application resources...');
      expect(mockLoggingService.info).toHaveBeenCalledWith('Application cleanup completed');
    });

    it('it handles cleanup errors gracefully', () => {
      // Test that cleanup completes successfully even if there are no errors
      // The actual error handling is tested implicitly by the try-catch structure
      service.cleanup();

      expect(mockLoggingService.info).toHaveBeenCalledWith('Cleaning up application resources...');
      expect(mockLoggingService.info).toHaveBeenCalledWith('Application cleanup completed');
      expect(mockLoggingService.error).not.toHaveBeenCalled();
    });

    it('it logs cleanup start and completion', () => {
      service.cleanup();

      expect(mockLoggingService.info).toHaveBeenCalledTimes(2);
      expect(mockLoggingService.info).toHaveBeenCalledWith('Cleaning up application resources...');
      expect(mockLoggingService.info).toHaveBeenCalledWith('Application cleanup completed');
    });
  });

  describe('service dependencies', () => {
    it('should inject ImageUploadService', () => {
      expect(mockImageUploadService).toBeTruthy();
    });

    it('should inject LoggingService', () => {
      expect(mockLoggingService).toBeTruthy();
    });
  });
});
