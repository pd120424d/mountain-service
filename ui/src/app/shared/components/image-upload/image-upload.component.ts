import { Component, Input, Output, EventEmitter, OnInit, OnDestroy, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';
import { ImageUploadService, ImageUploadProgress, CachedImage } from '../../../services/image-upload.service';
import { Subscription } from 'rxjs';

export interface ImageUploadEvent {
  file: File;
  preview: string;
  isValid: boolean;
  error?: string;
}

@Component({
  selector: 'app-image-upload',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <div class="image-upload-container">
      <!-- Upload Area -->
      <div
        class="upload-area"
        [class.drag-over]="isDragOver"
        [class.has-image]="previewUrl"
        (click)="triggerFileInput()"
        (dragover)="onDragOver($event)"
        (dragleave)="onDragLeave($event)"
        (drop)="onDrop($event)">

        <!-- File Input -->
        <input
          #fileInput
          type="file"
          accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
          (change)="onFileSelected($event)"
          style="display: none;">

        <!-- Preview Image -->
        <div *ngIf="previewUrl" class="image-preview">
          <img [src]="previewUrl" [alt]="'IMAGE_UPLOAD.PREVIEW_ALT' | translate">
          <div class="image-overlay">
            <button
              type="button"
              class="btn btn-sm btn-danger remove-btn"
              (click)="removeImage($event)"
              [title]="'IMAGE_UPLOAD.REMOVE' | translate">
              ✕
            </button>
          </div>
        </div>

        <!-- Upload Prompt -->
        <div *ngIf="!previewUrl" class="upload-prompt">
          <div class="upload-icon">📷</div>
          <p class="upload-text">{{ 'IMAGE_UPLOAD.DRAG_DROP_TEXT' | translate }}</p>
          <p class="upload-subtext">{{ 'IMAGE_UPLOAD.OR_CLICK_TEXT' | translate }}</p>
          <p class="upload-requirements">{{ 'IMAGE_UPLOAD.REQUIREMENTS' | translate }}</p>
        </div>

        <!-- Upload Progress -->
        <div *ngIf="uploadProgress && uploadProgress.status === 'uploading'" class="upload-progress">
          <div class="progress-bar">
            <div class="progress-fill" [style.width.%]="uploadProgress.progress"></div>
          </div>
          <p class="progress-text">{{ uploadProgress.message }}</p>
        </div>
      </div>

      <!-- Error Message -->
      <div *ngIf="errorMessage" class="error-message">
        <span class="error-icon">⚠️</span>
        {{ errorMessage }}
      </div>

      <!-- Success Message -->
      <div *ngIf="successMessage" class="success-message">
        <span class="success-icon">✅</span>
        {{ successMessage }}
      </div>
    </div>
  `,
  styles: [`
    .image-upload-container {
      width: 100%;
      max-width: 300px;
    }

    .upload-area {
      border: 2px dashed #ccc;
      border-radius: 8px;
      padding: 20px;
      text-align: center;
      cursor: pointer;
      transition: all 0.3s ease;
      position: relative;
      min-height: 200px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .upload-area:hover {
      border-color: #007bff;
      background-color: #f8f9fa;
    }

    .upload-area.drag-over {
      border-color: #007bff;
      background-color: #e3f2fd;
    }

    .upload-area.has-image {
      padding: 0;
      border: none;
    }

    .image-preview {
      position: relative;
      width: 100%;
      height: 200px;
      border-radius: 8px;
      overflow: hidden;
    }

    .image-preview img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.5);
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: opacity 0.3s ease;
    }

    .image-preview:hover .image-overlay {
      opacity: 1;
    }

    .remove-btn {
      background: #dc3545;
      border: none;
      color: white;
      width: 30px;
      height: 30px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 16px;
      cursor: pointer;
    }

    .upload-prompt {
      color: #666;
    }

    .upload-icon {
      font-size: 48px;
      margin-bottom: 16px;
    }

    .upload-text {
      font-size: 16px;
      font-weight: 500;
      margin-bottom: 8px;
    }

    .upload-subtext {
      font-size: 14px;
      margin-bottom: 12px;
    }

    .upload-requirements {
      font-size: 12px;
      color: #999;
    }

    .upload-progress {
      position: absolute;
      bottom: 10px;
      left: 10px;
      right: 10px;
    }

    .progress-bar {
      width: 100%;
      height: 4px;
      background: #e0e0e0;
      border-radius: 2px;
      overflow: hidden;
    }

    .progress-fill {
      height: 100%;
      background: #007bff;
      transition: width 0.3s ease;
    }

    .progress-text {
      font-size: 12px;
      margin-top: 4px;
      color: #666;
    }

    .error-message {
      margin-top: 8px;
      padding: 8px 12px;
      background: #f8d7da;
      color: #721c24;
      border: 1px solid #f5c6cb;
      border-radius: 4px;
      font-size: 14px;
    }

    .success-message {
      margin-top: 8px;
      padding: 8px 12px;
      background: #d4edda;
      color: #155724;
      border: 1px solid #c3e6cb;
      border-radius: 4px;
      font-size: 14px;
    }

    .error-icon, .success-icon {
      margin-right: 8px;
    }
  `]
})
export class ImageUploadComponent implements OnInit, OnDestroy, OnChanges {
  @Input() employeeId?: number;
  @Input() cacheKey?: string;
  @Input() initialImageUrl?: string | undefined;
  @Output() imageSelected = new EventEmitter<ImageUploadEvent>();
  @Output() imageUploaded = new EventEmitter<any>();
  @Output() imageRemoved = new EventEmitter<void>();

  previewUrl: string | null = null;
  selectedFile: File | null = null;
  isDragOver = false;
  uploadProgress: ImageUploadProgress | null = null;
  errorMessage: string | null = null;
  successMessage: string | null = null;

  private uploadSubscription?: Subscription;

  constructor(private imageUploadService: ImageUploadService) {}

  ngOnInit(): void {
    // Load cached image if available
    if (this.cacheKey) {
      const cached = this.imageUploadService.getCachedImage(this.cacheKey);
      if (cached) {
        this.previewUrl = cached.preview;
        this.selectedFile = cached.file;

        // If we have upload result, show success message
        if (cached.uploadResult) {
          this.successMessage = 'Image loaded from cache';
        }
      }
    }

    // Set initial image if provided and no cached image
    if (this.initialImageUrl && !this.previewUrl) {
      this.previewUrl = this.initialImageUrl;
    }
  }
  ngOnChanges(changes: SimpleChanges): void {
    if (changes['initialImageUrl'] && !this.previewUrl && this.initialImageUrl) {
      this.previewUrl = this.initialImageUrl;
    }
  }


  ngOnDestroy(): void {
    if (this.uploadSubscription) {
      this.uploadSubscription.unsubscribe();
    }
  }

  triggerFileInput(): void {
    // Don't trigger file input if an image is already selected
    if (this.previewUrl) {
      return;
    }

    const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement;
    fileInput?.click();
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      this.handleFile(input.files[0]);
    }
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    this.isDragOver = true;
  }

  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    this.isDragOver = false;
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    this.isDragOver = false;

    if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
      this.handleFile(event.dataTransfer.files[0]);
    }
  }

  removeImage(event: Event): void {
    event.stopPropagation();
    this.previewUrl = null;
    this.selectedFile = null;
    this.errorMessage = null;
    this.successMessage = null;
    this.uploadProgress = null;

    if (this.cacheKey) {
      this.imageUploadService.removeFromCache(this.cacheKey);
    }

    this.imageRemoved.emit();
  }

  private async handleFile(file: File): Promise<void> {
    this.clearMessages();

    // Validate file
    const validation = this.imageUploadService.validateImageFile(file);
    if (!validation.valid) {
      this.errorMessage = validation.error || 'Invalid file';
      this.imageSelected.emit({
        file,
        preview: '',
        isValid: false,
        error: validation.error
      });
      return;
    }

    try {
      // Create preview
      const preview = await this.imageUploadService.createImagePreview(file);
      this.previewUrl = preview;
      this.selectedFile = file;

      // Cache the image
      if (this.cacheKey) {
        this.imageUploadService.cacheImage(this.cacheKey, file, preview);
      }

      // Emit event
      this.imageSelected.emit({
        file,
        preview,
        isValid: true
      });

    } catch (error) {
      this.errorMessage = 'Failed to process image';
    }
  }

  private clearMessages(): void {
    this.errorMessage = null;
    this.successMessage = null;
  }
}
