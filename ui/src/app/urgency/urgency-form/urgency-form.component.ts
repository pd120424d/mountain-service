import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { UrgencyService } from '../urgency.service';
import { UrgencyCreateRequest, UrgencyLevel, EnhancedLocation, LocationCoordinates, formatLocationForApi, parseLocationString } from '../../shared/models';
import { LocationMapComponent } from '../../shared/components';

@Component({
  selector: 'app-urgency-form',
  standalone: true,
  imports: [RouterModule, TranslateModule, CommonModule, ReactiveFormsModule, LocationMapComponent],
  templateUrl: './urgency-form.component.html',
  styleUrls: ['./urgency-form.component.css']
})
export class UrgencyFormComponent extends BaseTranslatableComponent implements OnInit {
  urgencyForm!: FormGroup;
  urgencyLevels = Object.values(UrgencyLevel);
  isSubmitting = false;

  currentLocation?: EnhancedLocation;
  showMap = false;

  constructor(
    private fb: FormBuilder,
    private urgencyService: UrgencyService,
    private router: Router,
    private toastr: ToastrService,
    translate: TranslateService
  ) {
    super(translate);
  }

  ngOnInit(): void {
    this.initializeForm();
  }

  private initializeForm(): void {
    this.urgencyForm = this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      contactPhone: ['', [Validators.required, Validators.minLength(8)]],
      location: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      level: [UrgencyLevel.Medium, [Validators.required]]
    });
  }

  onSubmit(): void {
    if (this.urgencyForm.valid && !this.isSubmitting) {
      this.isSubmitting = true;

      const formValue = this.urgencyForm.value;
      const urgencyRequest: UrgencyCreateRequest = {
        ...formValue,
        location: this.currentLocation ? formatLocationForApi(this.currentLocation) : formValue.location
      };

      this.urgencyService.addUrgency(urgencyRequest).subscribe({
        next: (response) => {
          console.log('Urgency created successfully:', response);
          this.toastr.success(this.translate.instant('URGENCY_FORM.SUCCESS_MESSAGE'));
          this.router.navigate(['/']);
        },
        error: (error) => {
          console.error('Error submitting urgency:', error);
          this.isSubmitting = false;
        }
      });
    } else {
      this.markFormGroupTouched();
    }
  }

  cancel(): void {
    this.router.navigate(['/']);
  }

  private markFormGroupTouched(): void {
    Object.keys(this.urgencyForm.controls).forEach(key => {
      const control = this.urgencyForm.get(key);
      control?.markAsTouched();
    });
  }

  // Helper methods for template
  isFieldInvalid(fieldName: string): boolean {
    const field = this.urgencyForm.get(fieldName);
    return !!(field && field.invalid && (field.dirty || field.touched));
  }

  getFieldError(fieldName: string): string {
    const field = this.urgencyForm.get(fieldName);
    if (field?.errors) {
      if (field.errors['required']) return `${fieldName.toUpperCase()}_FORM.REQUIRED`;
      if (field.errors['email']) return `${fieldName.toUpperCase()}_FORM.INVALID_EMAIL`;
      if (field.errors['minlength']) return `${fieldName.toUpperCase()}_FORM.MIN_LENGTH`;
    }
    return '';
  }

  toggleMap(): void {
    this.showMap = !this.showMap;

    if (this.showMap && !this.currentLocation) {
      const locationText = this.urgencyForm.get('location')?.value;
      if (locationText) {
        this.currentLocation = parseLocationString(locationText) || {
          text: locationText,
          source: 'manual'
        };
      }
    }
  }

  onLocationSelected(location: EnhancedLocation): void {
    this.currentLocation = location;

    this.urgencyForm.patchValue({
      location: location.text
    });

    const locationControl = this.urgencyForm.get('location');
    if (locationControl) {
      locationControl.markAsTouched();
      locationControl.updateValueAndValidity();
    }
  }

  onCoordinatesChanged(coordinates: LocationCoordinates): void {
    if (this.currentLocation) {
      this.currentLocation = {
        ...this.currentLocation,
        coordinates,
        source: 'map'
      };
    }
  }

  clearMapLocation(): void {
    this.currentLocation = undefined;
    this.urgencyForm.patchValue({
      location: ''
    });
  }

  hasCoordinates(): boolean {
    return !!(this.currentLocation?.coordinates);
  }
}