import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ImageUploadService, ImageUploadResult } from './image-upload.service';
import { LoggingService } from './logging.service';
import { SessionStorageService } from './session-storage.service';
import { environment } from '../../environments/environment';

describe('ImageUploadService', () => {
  let service: ImageUploadService;
  let httpMock: HttpTestingController;
  let mockLoggingService: jasmine.SpyObj<LoggingService>;
  let mockSessionStorageService: jasmine.SpyObj<SessionStorageService>;

  beforeEach(() => {
    const loggingSpy = jasmine.createSpyObj('LoggingService', ['info', 'error', 'debug', 'warn']);
    const sessionStorageSpy = jasmine.createSpyObj('SessionStorageService', [
      'storeImageData', 'getImageData', 'removeImageData', 'clearAllImageData', 
      'getAllImageKeys', 'getStorageInfo', 'getImagePreferences', 'storeImagePreferences'
    ]);

    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [
        ImageUploadService,
        { provide: LoggingService, useValue: loggingSpy },
        { provide: SessionStorageService, useValue: sessionStorageSpy }
      ]
    });

    service = TestBed.inject(ImageUploadService);
    httpMock = TestBed.inject(HttpTestingController);
    mockLoggingService = TestBed.inject(LoggingService) as jasmine.SpyObj<LoggingService>;
    mockSessionStorageService = TestBed.inject(SessionStorageService) as jasmine.SpyObj<SessionStorageService>;
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe('uploadProfilePicture', () => {
    it('should upload profile picture successfully', (done) => {
      const employeeId = 123;
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      const expectedResult: ImageUploadResult = {
        blobUrl: 'https://test.blob.core.windows.net/container/test.jpg',
        blobName: 'employee-123/test.jpg',
        size: 1024,
        message: 'Upload completed successfully'
      };

      service.uploadProfilePicture(employeeId, file).subscribe({
        next: (progress) => {
          if (progress.status === 'completed') {
            expect(progress.progress).toBe(100);
            expect(progress.message).toBe('Upload completed successfully');
            expect(mockLoggingService.info).toHaveBeenCalledWith(
              `Profile picture uploaded successfully for employee ${employeeId}`
            );
            done();
          }
        },
        error: () => fail('Should not have errored')
      });

      const req = httpMock.expectOne(`${environment.apiUrl}/employees/${employeeId}/profile-picture`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toBeInstanceOf(FormData);

      // Simulate successful response
      req.event({
        type: 4, // HttpEventType.Response
        body: expectedResult
      } as any);
    });

    xit('should handle upload error', (done) => {
      const employeeId = 123;
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });

      service.uploadProfilePicture(employeeId, file).subscribe({
        next: () => fail('Should not have succeeded'),
        error: (error) => {
          expect(error.status).toBe('error');
          expect(error.message).toContain('Http failure response');
          expect(mockLoggingService.error).toHaveBeenCalledWith(
            `Failed to upload profile picture for employee ${employeeId}:`,
            jasmine.any(Object)
          );
          done();
        }
      });

      const req = httpMock.expectOne(`${environment.apiUrl}/employees/${employeeId}/profile-picture`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toBeInstanceOf(FormData);

      // Simulate error response
      req.error(new ErrorEvent('Network error'), { status: 500, statusText: 'Internal Server Error' });
    });
  });

  describe('deleteProfilePicture', () => {
    it('should delete profile picture successfully', () => {
      const employeeId = 123;
      const blobName = 'employee-123/test.jpg';

      service.deleteProfilePicture(employeeId, blobName).subscribe({
        next: (response) => {
          expect(response).toBeDefined();
          expect(mockLoggingService.info).toHaveBeenCalledWith(
            `Profile picture deleted successfully for employee ${employeeId}`
          );
          expect(mockSessionStorageService.removeImageData).toHaveBeenCalledWith(`employee-${employeeId}`);
        },
        error: () => fail('Should not have errored')
      });

      const req = httpMock.expectOne(
        `${environment.apiUrl}/employees/${employeeId}/profile-picture?blobName=${encodeURIComponent(blobName)}`
      );
      expect(req.request.method).toBe('DELETE');
      req.flush({ message: 'Deleted successfully' });
    });

    it('should handle delete error', () => {
      const employeeId = 123;
      const blobName = 'employee-123/test.jpg';

      service.deleteProfilePicture(employeeId, blobName).subscribe({
        next: () => fail('Should not have succeeded'),
        error: (error) => {
          expect(error).toBeDefined();
          expect(mockLoggingService.error).toHaveBeenCalledWith(
            `Failed to delete profile picture for employee ${employeeId}:`,
            jasmine.any(Object)
          );
        }
      });

      const req = httpMock.expectOne(
        `${environment.apiUrl}/employees/${employeeId}/profile-picture?blobName=${encodeURIComponent(blobName)}`
      );
      req.flush({ error: 'Delete failed' }, { status: 500, statusText: 'Internal Server Error' });
    });
  });

  describe('validateImageFile', () => {
    it('should validate valid JPEG file', () => {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      Object.defineProperty(file, 'size', { value: 1024 * 1024 }); // 1MB

      const result = service.validateImageFile(file);

      expect(result.valid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    it('should validate valid PNG file', () => {
      const file = new File(['test'], 'test.png', { type: 'image/png' });
      Object.defineProperty(file, 'size', { value: 2 * 1024 * 1024 }); // 2MB

      const result = service.validateImageFile(file);

      expect(result.valid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    it('should reject file that is too large', () => {
      const file = new File(['test'], 'large.jpg', { type: 'image/jpeg' });
      Object.defineProperty(file, 'size', { value: 6 * 1024 * 1024 }); // 6MB

      const result = service.validateImageFile(file);

      expect(result.valid).toBe(false);
      expect(result.error).toBe('File size exceeds 5MB limit');
    });

    it('should reject invalid file type', () => {
      const file = new File(['test'], 'test.txt', { type: 'text/plain' });
      Object.defineProperty(file, 'size', { value: 1024 });

      const result = service.validateImageFile(file);

      expect(result.valid).toBe(false);
      expect(result.error).toBe('Invalid file type. Please select a JPEG, PNG, GIF, or WebP image.');
    });
  });

  describe('createImagePreview', () => {
    it('should create image preview successfully', async () => {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      
      // Mock FileReader
      const mockFileReader = {
        readAsDataURL: jasmine.createSpy('readAsDataURL').and.callFake(function(this: any) {
          setTimeout(() => {
            this.onload({ target: { result: 'data:image/jpeg;base64,test' } });
          }, 0);
        }),
        onload: null,
        onerror: null
      };
      
      spyOn(window, 'FileReader').and.returnValue(mockFileReader as any);

      const preview = await service.createImagePreview(file);

      expect(preview).toBe('data:image/jpeg;base64,test');
      expect(mockFileReader.readAsDataURL).toHaveBeenCalledWith(file);
    });

    it('should handle file read error', async () => {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      
      // Mock FileReader with error
      const mockFileReader = {
        readAsDataURL: jasmine.createSpy('readAsDataURL').and.callFake(function(this: any) {
          setTimeout(() => {
            this.onerror();
          }, 0);
        }),
        onload: null,
        onerror: null
      };
      
      spyOn(window, 'FileReader').and.returnValue(mockFileReader as any);

      try {
        await service.createImagePreview(file);
        fail('Should have thrown an error');
      } catch (error) {
        expect(error).toEqual(new Error('Failed to read file'));
      }
    });
  });

  describe('cacheImage', () => {
    it('should cache image in memory and session storage', () => {
      const key = 'test-key';
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      const preview = 'data:image/jpeg;base64,test';
      const uploadResult: ImageUploadResult = {
        blobUrl: 'https://test.blob.core.windows.net/test.jpg',
        blobName: 'test.jpg',
        size: 1024,
        message: 'Success'
      };

      service.cacheImage(key, file, preview, uploadResult);

      expect(mockSessionStorageService.storeImageData).toHaveBeenCalledWith(key, {
        preview,
        fileName: file.name,
        fileSize: file.size,
        fileType: file.type,
        uploadResult,
        timestamp: jasmine.any(Number)
      });
      expect(mockLoggingService.debug).toHaveBeenCalledWith(`Image cached with key: ${key}`);
    });
  });

  describe('getCachedImage', () => {
    it('should return cached image from memory', () => {
      const key = 'test-key';
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      const preview = 'data:image/jpeg;base64,test';

      // First cache the image
      service.cacheImage(key, file, preview);

      const cached = service.getCachedImage(key);

      expect(cached).toBeDefined();
      expect(cached!.file).toBe(file);
      expect(cached!.preview).toBe(preview);
    });

    it('should return cached image from session storage when not in memory', () => {
      const key = 'test-key';
      const sessionData = {
        preview: 'data:image/jpeg;base64,test',
        fileName: 'test.jpg',
        fileSize: 1024,
        fileType: 'image/jpeg',
        timestamp: Date.now()
      };

      mockSessionStorageService.getImageData.and.returnValue(sessionData);

      const cached = service.getCachedImage(key);

      expect(cached).toBeDefined();
      expect(cached!.preview).toBe(sessionData.preview);
      expect(cached!.file.name).toBe(sessionData.fileName);
      expect(cached!.file.type).toBe(sessionData.fileType);
      expect(mockSessionStorageService.getImageData).toHaveBeenCalledWith(key);
    });

    it('should return null when no cached image exists', () => {
      const key = 'non-existent-key';
      mockSessionStorageService.getImageData.and.returnValue(null);

      const cached = service.getCachedImage(key);

      expect(cached).toBeNull();
    });
  });

  describe('clearCache', () => {
    it('should clear both memory and session storage cache', () => {
      service.clearCache();

      expect(mockSessionStorageService.clearAllImageData).toHaveBeenCalled();
      expect(mockLoggingService.info).toHaveBeenCalledWith('Image cache cleared');
    });
  });
});
