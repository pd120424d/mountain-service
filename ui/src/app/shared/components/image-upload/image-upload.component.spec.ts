import { ComponentFixture, TestBed } from '@angular/core/testing';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ImageUploadComponent, ImageUploadEvent } from './image-upload.component';
import { ImageUploadService, CachedImage } from '../../../services/image-upload.service';

describe('ImageUploadComponent', () => {
  let component: ImageUploadComponent;
  let fixture: ComponentFixture<ImageUploadComponent>;
  let mockImageUploadService: jasmine.SpyObj<ImageUploadService>;

  beforeEach(async () => {
    const imageUploadSpy = jasmine.createSpyObj('ImageUploadService', [
      'validateImageFile', 'createImagePreview', 'getCachedImage', 'cacheImage', 'removeFromCache'
    ]);

    await TestBed.configureTestingModule({
      imports: [ImageUploadComponent, TranslateModule.forRoot()],
      providers: [
        { provide: ImageUploadService, useValue: imageUploadSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ImageUploadComponent);
    component = fixture.componentInstance;
    mockImageUploadService = TestBed.inject(ImageUploadService) as jasmine.SpyObj<ImageUploadService>;
  });

  describe('ngOnInit', () => {
    it('should load cached image when cache key is provided', () => {
      const cachedImage: CachedImage = {
        file: new File(['test'], 'test.jpg', { type: 'image/jpeg' }),
        preview: 'data:image/jpeg;base64,test',
        timestamp: Date.now()
      };

      component.cacheKey = 'test-key';
      mockImageUploadService.getCachedImage.and.returnValue(cachedImage);

      component.ngOnInit();

      expect(mockImageUploadService.getCachedImage).toHaveBeenCalledWith('test-key');
      expect(component.previewUrl).toBe(cachedImage.preview);
      expect(component.selectedFile).toBe(cachedImage.file);
    });

    it('should set initial image URL when no cached image exists', () => {
      const initialUrl = 'https://example.com/image.jpg';
      component.initialImageUrl = initialUrl;
      component.cacheKey = 'test-key';
      mockImageUploadService.getCachedImage.and.returnValue(null);

      component.ngOnInit();

      expect(component.previewUrl).toBe(initialUrl);
    });

    it('should show success message when cached image has upload result', () => {
      const cachedImage: CachedImage = {
        file: new File(['test'], 'test.jpg', { type: 'image/jpeg' }),
        preview: 'data:image/jpeg;base64,test',
        uploadResult: {
          blobUrl: 'https://test.blob.core.windows.net/test.jpg',
          blobName: 'test.jpg',
          size: 1024,
          message: 'Success'
        },
        timestamp: Date.now()
      };

      component.cacheKey = 'test-key';
      mockImageUploadService.getCachedImage.and.returnValue(cachedImage);

      component.ngOnInit();

      expect(component.successMessage).toBe('Image loaded from cache');
    });
  });

  describe('handleFile', () => {
    it('should handle valid image file successfully', async () => {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      const preview = 'data:image/jpeg;base64,test';

      mockImageUploadService.validateImageFile.and.returnValue({ valid: true });
      mockImageUploadService.createImagePreview.and.returnValue(Promise.resolve(preview));
      component.cacheKey = 'test-key';

      spyOn(component.imageSelected, 'emit');

      await component['handleFile'](file);

      expect(mockImageUploadService.validateImageFile).toHaveBeenCalledWith(file);
      expect(mockImageUploadService.createImagePreview).toHaveBeenCalledWith(file);
      expect(mockImageUploadService.cacheImage).toHaveBeenCalledWith('test-key', file, preview);
      expect(component.previewUrl).toBe(preview);
      expect(component.selectedFile).toBe(file);
      expect(component.imageSelected.emit).toHaveBeenCalledWith({
        file,
        preview,
        isValid: true
      });
    });

    it('should handle invalid image file', async () => {
      const file = new File(['test'], 'test.txt', { type: 'text/plain' });
      const errorMessage = 'Invalid file type';

      mockImageUploadService.validateImageFile.and.returnValue({ 
        valid: false, 
        error: errorMessage 
      });

      spyOn(component.imageSelected, 'emit');

      await component['handleFile'](file);

      expect(component.errorMessage).toBe(errorMessage);
      expect(component.imageSelected.emit).toHaveBeenCalledWith({
        file,
        preview: '',
        isValid: false,
        error: errorMessage
      });
    });

    it('should handle image preview creation error', async () => {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });

      mockImageUploadService.validateImageFile.and.returnValue({ valid: true });
      mockImageUploadService.createImagePreview.and.returnValue(Promise.reject(new Error('Preview failed')));

      await component['handleFile'](file);

      expect(component.errorMessage).toBe('Failed to process image');
    });
  });

  describe('onFileSelected', () => {
    it('should handle file selection', async () => {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      const mockEvent = {
        target: {
          files: [file]
        }
      } as any;

      mockImageUploadService.validateImageFile.and.returnValue({ valid: true });
      mockImageUploadService.createImagePreview.and.returnValue(Promise.resolve('data:image/jpeg;base64,test'));

      spyOn(component, 'handleFile' as any);

      component.onFileSelected(mockEvent);

      expect(component['handleFile']).toHaveBeenCalledWith(file);
    });

    it('should do nothing when no file is selected', () => {
      const mockEvent = {
        target: {
          files: null
        }
      } as any;

      spyOn(component, 'handleFile' as any);

      component.onFileSelected(mockEvent);

      expect(component['handleFile']).not.toHaveBeenCalled();
    });
  });

  describe('onDrop', () => {
    it('should handle file drop', async () => {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      const mockEvent = {
        preventDefault: jasmine.createSpy('preventDefault'),
        dataTransfer: {
          files: [file]
        }
      } as any;

      mockImageUploadService.validateImageFile.and.returnValue({ valid: true });
      mockImageUploadService.createImagePreview.and.returnValue(Promise.resolve('data:image/jpeg;base64,test'));

      spyOn(component, 'handleFile' as any);

      component.onDrop(mockEvent);

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(component.isDragOver).toBe(false);
      expect(component['handleFile']).toHaveBeenCalledWith(file);
    });

    it('should do nothing when no file is dropped', () => {
      const mockEvent = {
        preventDefault: jasmine.createSpy('preventDefault'),
        dataTransfer: {
          files: null
        }
      } as any;

      spyOn(component, 'handleFile' as any);

      component.onDrop(mockEvent);

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(component['handleFile']).not.toHaveBeenCalled();
    });
  });

  describe('removeImage', () => {
    it('should remove image and clear state', () => {
      const mockEvent = {
        stopPropagation: jasmine.createSpy('stopPropagation')
      } as any;

      component.previewUrl = 'test-preview';
      component.selectedFile = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      component.errorMessage = 'test error';
      component.successMessage = 'test success';
      component.cacheKey = 'test-key';

      spyOn(component.imageRemoved, 'emit');

      component.removeImage(mockEvent);

      expect(mockEvent.stopPropagation).toHaveBeenCalled();
      expect(component.previewUrl).toBeNull();
      expect(component.selectedFile).toBeNull();
      expect(component.errorMessage).toBeNull();
      expect(component.successMessage).toBeNull();
      expect(mockImageUploadService.removeFromCache).toHaveBeenCalledWith('test-key');
      expect(component.imageRemoved.emit).toHaveBeenCalled();
    });
  });

  describe('onDragOver', () => {
    it('should set drag over state', () => {
      const mockEvent = {
        preventDefault: jasmine.createSpy('preventDefault')
      } as any;

      component.onDragOver(mockEvent);

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(component.isDragOver).toBe(true);
    });
  });

  describe('onDragLeave', () => {
    it('should clear drag over state', () => {
      const mockEvent = {
        preventDefault: jasmine.createSpy('preventDefault')
      } as any;

      component.isDragOver = true;

      component.onDragLeave(mockEvent);

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(component.isDragOver).toBe(false);
    });
  });

  describe('triggerFileInput', () => {
    it('should not trigger file input when image is already selected', () => {
      component.previewUrl = 'test-preview';
      
      spyOn(document, 'querySelector');

      component.triggerFileInput();

      expect(document.querySelector).not.toHaveBeenCalled();
    });

    it('should trigger file input when no image is selected', () => {
      component.previewUrl = null;
      
      const mockInput = {
        click: jasmine.createSpy('click')
      };
      spyOn(document, 'querySelector').and.returnValue(mockInput as any);

      component.triggerFileInput();

      expect(document.querySelector).toHaveBeenCalledWith('input[type="file"]');
      expect(mockInput.click).toHaveBeenCalled();
    });
  });
});
