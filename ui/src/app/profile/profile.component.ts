import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { NgxSpinnerService } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';

import { AuthService } from '../services/auth.service';
import { EmployeeService } from '../employee/employee.service';
import { ImageUploadService } from '../services/image-upload.service';
import { ImageUploadEvent } from '../shared/components/image-upload/image-upload.component';
import { Employee, EmployeeUpdateRequest } from '../shared/models';
import { BaseTranslatableComponent } from '../base-translatable.component';
import { ImageUploadComponent } from '../shared/components/image-upload/image-upload.component';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule, ImageUploadComponent],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent extends BaseTranslatableComponent implements OnInit {
  profileForm: FormGroup;
  currentUser: Employee | null = null;
  isSubmitting = false;
  selectedImageFile: File | null = null;
  currentProfilePictureUrl?: string;

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private employeeService: EmployeeService,
    private imageUploadService: ImageUploadService,
    private router: Router,
    private spinner: NgxSpinnerService,
    private toastr: ToastrService,
    translate: TranslateService
  ) {
    super(translate);
    
    this.profileForm = this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      phone: ['', [Validators.required]],
      username: ['', [Validators.required, Validators.minLength(3)]],
      profilePicture: [null]
    });
  }

  ngOnInit(): void {
    this.loadCurrentUser();
  }

  private loadCurrentUser(): void {
    const userId = this.authService.getUserId();
    if (!userId) {
      this.router.navigate(['/login']);
      return;
    }

    this.spinner.show();
    this.employeeService.getEmployeeById(parseInt(userId)).subscribe({
      next: (user) => {
        this.currentUser = user;
        this.currentProfilePictureUrl = user.profilePicture;
        this.profileForm.patchValue({
          firstName: user.firstName,
          lastName: user.lastName,
          email: user.email,
          phone: user.phone,
          username: user.username
        });
        this.spinner.hide();
      },
      error: (error) => {
        this.spinner.hide();
        this.toastr.error(this.translate.instant('PROFILE.LOAD_ERROR'));
        console.error('Error loading user profile:', error);
      }
    });
  }

  onImageSelected(event: ImageUploadEvent): void {
    if (event.isValid && event.file) {
      this.selectedImageFile = event.file;
      this.profileForm.patchValue({ profilePicture: event.preview });
    } else {
      this.selectedImageFile = null;
      this.toastr.error(event.error || 'Invalid image file');
    }
  }

  onImageRemoved(): void {
    this.selectedImageFile = null;
    this.currentProfilePictureUrl = undefined;
    this.profileForm.patchValue({ profilePicture: null });
  }

  onSubmit(): void {
    if (this.profileForm.invalid || !this.currentUser) {
      this.markFormGroupTouched();
      return;
    }

    this.isSubmitting = true;
    this.spinner.show();

    const updateRequest: EmployeeUpdateRequest = {
      firstName: this.profileForm.value.firstName,
      lastName: this.profileForm.value.lastName,
      email: this.profileForm.value.email,
      phone: this.profileForm.value.phone,
      username: this.profileForm.value.username,
      profileType: this.currentUser.profileType,
      gender: this.currentUser.gender
    };

    this.employeeService.updateEmployee(this.currentUser.id!, updateRequest).subscribe({
      next: (updatedUser) => {
        if (this.selectedImageFile) {
          this.uploadImageAfterUpdate(updatedUser.id!);
        } else {
          this.spinner.hide();
          this.isSubmitting = false;
          this.toastr.success(this.translate.instant('PROFILE.UPDATE_SUCCESS'));
        }
      },
      error: (error) => {
        this.spinner.hide();
        this.isSubmitting = false;
        this.toastr.error(this.translate.instant('PROFILE.UPDATE_ERROR'));
        console.error('Error updating profile:', error);
      }
    });
  }

  private uploadImageAfterUpdate(employeeId: number): void {
    if (!this.selectedImageFile) {
      this.spinner.hide();
      this.isSubmitting = false;
      this.toastr.success(this.translate.instant('PROFILE.UPDATE_SUCCESS'));
      return;
    }

    this.imageUploadService.uploadProfilePicture(employeeId, this.selectedImageFile).subscribe({
      next: (progress) => {
        if (progress.status === 'completed') {
          this.spinner.hide();
          this.isSubmitting = false;
          this.toastr.success(this.translate.instant('PROFILE.UPDATE_SUCCESS'));
          this.loadCurrentUser(); // Reload to get updated profile picture
        }
      },
      error: (error) => {
        this.spinner.hide();
        this.isSubmitting = false;
        this.toastr.warning('Profile updated but image upload failed: ' + error.message);
      }
    });
  }

  private markFormGroupTouched(): void {
    Object.keys(this.profileForm.controls).forEach(key => {
      const control = this.profileForm.get(key);
      control?.markAsTouched();
    });
  }

  goBack(): void {
    this.router.navigate(['/']);
  }
}
