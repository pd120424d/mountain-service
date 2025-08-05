import { TestBed } from '@angular/core/testing';
import { SessionStorageService, SessionImageData } from './session-storage.service';
import { LoggingService } from './logging.service';

describe('SessionStorageService', () => {
  let service: SessionStorageService;
  let mockLoggingService: jasmine.SpyObj<LoggingService>;

  const mockImageData: SessionImageData = {
    preview: 'data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD',
    fileName: 'test-image.jpg',
    fileSize: 1024,
    fileType: 'image/jpeg',
    uploadResult: {
      blobUrl: 'https://example.com/blob/test-image.jpg',
      blobName: 'test-image-blob.jpg',
      size: 1024,
      message: 'Upload successful'
    },
    timestamp: Date.now()
  };

  beforeEach(() => {
    const loggingServiceSpy = jasmine.createSpyObj('LoggingService', ['info', 'debug', 'error']);

    TestBed.configureTestingModule({
      providers: [
        SessionStorageService,
        { provide: LoggingService, useValue: loggingServiceSpy }
      ]
    });

    service = TestBed.inject(SessionStorageService);
    mockLoggingService = TestBed.inject(LoggingService) as jasmine.SpyObj<LoggingService>;

    // Clear sessionStorage before each test
    sessionStorage.clear();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should initialize and log initialization', () => {
    expect(mockLoggingService.info).toHaveBeenCalledWith('SessionStorageService initialized');
  });

  describe('storeImageData', () => {
    it('it succeeds when storing image data successfully', () => {
      service.storeImageData('test-key', mockImageData);

      const storedData = sessionStorage.getItem('mountain_service_image_test-key');
      expect(storedData).toBeTruthy();
      expect(mockLoggingService.debug).toHaveBeenCalledWith('Image data stored in session storage: test-key');
    });

    it('it handles storage errors gracefully', () => {
      spyOn(sessionStorage, 'setItem').and.throwError('Storage full');
      spyOn(service, 'cleanupExpiredImages' as any);

      service.storeImageData('test-key', mockImageData);

      expect(mockLoggingService.error).toHaveBeenCalledWith('Failed to store image data in session storage:', jasmine.any(Error));
    });
  });

  describe('getImageData', () => {
    it('it succeeds when retrieving valid image data', () => {
      const testData = { ...mockImageData, timestamp: Date.now() };
      sessionStorage.setItem('mountain_service_image_test-key', JSON.stringify(testData));

      const result = service.getImageData('test-key');

      expect(result).toEqual(testData);
    });

    it('it returns null when no data exists', () => {
      const result = service.getImageData('non-existent-key');

      expect(result).toBeNull();
    });

    it('it returns null and removes expired data', () => {
      const expiredData = { ...mockImageData, timestamp: Date.now() - (25 * 60 * 60 * 1000) }; // 25 hours ago
      sessionStorage.setItem('mountain_service_image_test-key', JSON.stringify(expiredData));
      spyOn(service, 'removeImageData');

      const result = service.getImageData('test-key');

      expect(result).toBeNull();
      expect(service.removeImageData).toHaveBeenCalledWith('test-key');
    });

    it('it handles JSON parse errors', () => {
      sessionStorage.setItem('mountain_service_image_test-key', 'invalid-json');

      const result = service.getImageData('test-key');

      expect(result).toBeNull();
      expect(mockLoggingService.error).toHaveBeenCalledWith('Failed to retrieve image data from session storage:', jasmine.any(Error));
    });
  });

  describe('removeImageData', () => {
    it('it succeeds when removing image data', () => {
      sessionStorage.setItem('mountain_service_image_test-key', JSON.stringify(mockImageData));

      service.removeImageData('test-key');

      expect(sessionStorage.getItem('mountain_service_image_test-key')).toBeNull();
      expect(mockLoggingService.debug).toHaveBeenCalledWith('Image data removed from session storage: test-key');
    });

    it('it handles removal errors', () => {
      spyOn(sessionStorage, 'removeItem').and.throwError('Remove error');

      service.removeImageData('test-key');

      expect(mockLoggingService.error).toHaveBeenCalledWith('Failed to remove image data from session storage:', jasmine.any(Error));
    });
  });

  describe('clearAllImageData', () => {
    it('it succeeds when clearing all image data', () => {
      sessionStorage.setItem('mountain_service_image_key1', JSON.stringify(mockImageData));
      sessionStorage.setItem('mountain_service_image_key2', JSON.stringify(mockImageData));
      sessionStorage.setItem('other_key', 'other_value');

      service.clearAllImageData();

      expect(sessionStorage.getItem('mountain_service_image_key1')).toBeNull();
      expect(sessionStorage.getItem('mountain_service_image_key2')).toBeNull();
      expect(sessionStorage.getItem('other_key')).toBe('other_value'); // Should not be removed
      expect(mockLoggingService.info).toHaveBeenCalledWith('Cleared 2 image data entries from session storage');
    });

    it('it handles clear errors', () => {
      spyOn(sessionStorage, 'removeItem').and.throwError('Clear error');
      sessionStorage.setItem('mountain_service_image_key1', JSON.stringify(mockImageData));

      service.clearAllImageData();

      expect(mockLoggingService.error).toHaveBeenCalledWith('Failed to clear image data from session storage:', jasmine.any(Error));
    });
  });

  describe('getAllImageKeys', () => {
    it('it succeeds when getting all image keys', () => {
      sessionStorage.setItem('mountain_service_image_key1', JSON.stringify(mockImageData));
      sessionStorage.setItem('mountain_service_image_key2', JSON.stringify(mockImageData));
      sessionStorage.setItem('other_key', 'other_value');

      const keys = service.getAllImageKeys();

      expect(keys.length).toBe(2);
      expect(keys).toContain('key1');
      expect(keys).toContain('key2');
    });

    it('it returns empty array when no image keys exist', () => {
      // Clear any existing keys
      sessionStorage.clear();

      const keys = service.getAllImageKeys();

      expect(keys).toEqual([]);
    });
  });

  describe('getStorageInfo', () => {
    it('it succeeds when getting storage information', () => {
      sessionStorage.setItem('mountain_service_image_key1', JSON.stringify(mockImageData));
      sessionStorage.setItem('mountain_service_image_key2', JSON.stringify(mockImageData));

      const info = service.getStorageInfo();

      expect(info.totalKeys).toBe(2);
      expect(info.estimatedSize).toContain('B');
      expect(info.availableSpace).toContain('B');
    });

    it('it returns default values when storage is empty', () => {
      // Clear any existing keys
      sessionStorage.clear();

      const info = service.getStorageInfo();

      expect(info.totalKeys).toBe(0);
      expect(info.estimatedSize).toBe('0 B');
      expect(info.availableSpace).toBe('5 MB'); // The service returns '5 MB' as default
    });
  });

  describe('storeImagePreferences', () => {
    it('it succeeds when storing image preferences', () => {
      const preferences = { autoUpload: false, compressionQuality: 0.9, maxFileSize: 10485760 };

      service.storeImagePreferences(preferences);

      const storedPrefs = sessionStorage.getItem('mountain_service_image_prefs');
      expect(storedPrefs).toBe(JSON.stringify(preferences));
      expect(mockLoggingService.debug).toHaveBeenCalledWith('Image preferences stored');
    });

    it('it handles storage errors when storing preferences', () => {
      spyOn(sessionStorage, 'setItem').and.throwError('Storage error');

      service.storeImagePreferences({ autoUpload: false });

      expect(mockLoggingService.error).toHaveBeenCalledWith('Failed to store image preferences:', jasmine.any(Error));
    });
  });

  describe('getImagePreferences', () => {
    it('it succeeds when getting stored preferences', () => {
      const preferences = { autoUpload: false, compressionQuality: 0.9, maxFileSize: 10485760 };
      sessionStorage.setItem('mountain_service_image_prefs', JSON.stringify(preferences));

      const result = service.getImagePreferences();

      expect(result).toEqual(preferences);
    });

    it('it returns defaults when no preferences stored', () => {
      const result = service.getImagePreferences();

      expect(result).toEqual({
        autoUpload: true,
        compressionQuality: 0.8,
        maxFileSize: 5 * 1024 * 1024
      });
    });

    it('it returns defaults with partial preferences', () => {
      sessionStorage.setItem('mountain_service_image_prefs', JSON.stringify({ autoUpload: false }));

      const result = service.getImagePreferences();

      expect(result.autoUpload).toBe(false);
      expect(result.compressionQuality).toBe(0.8);
      expect(result.maxFileSize).toBe(5 * 1024 * 1024);
    });

    it('it handles JSON parse errors and returns defaults', () => {
      sessionStorage.setItem('mountain_service_image_prefs', 'invalid-json');

      const result = service.getImagePreferences();

      expect(result).toEqual({
        autoUpload: true,
        compressionQuality: 0.8,
        maxFileSize: 5 * 1024 * 1024
      });
      expect(mockLoggingService.error).toHaveBeenCalledWith('Failed to get image preferences:', jasmine.any(Error));
    });
  });
});
