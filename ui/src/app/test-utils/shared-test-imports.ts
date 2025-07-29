import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { importProvidersFrom } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { ActivatedRoute } from '@angular/router';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ToastrModule, ToastrService } from 'ngx-toastr';
import { NgxSpinnerService } from 'ngx-spinner';
import { EMPTY } from 'rxjs';

export const sharedTestingProviders = [
  provideHttpClient(withInterceptorsFromDi()),
  provideHttpClientTesting(),

  // You can wrap TranslateModule in `importProvidersFrom` like this:
  importProvidersFrom(TranslateModule.forRoot()),
  importProvidersFrom(ToastrModule.forRoot()),
  {
    provide: ActivatedRoute,
    useValue: {
      snapshot: {
        params: {},
        queryParams: {},
        data: {}
      }
    }
  },

  // ✅ Dialog mocks
  { provide: MatDialogRef, useValue: {} },
  { provide: MAT_DIALOG_DATA, useValue: {} },

  // ✅ NgxSpinner and Toastr service mocks
  {
    provide: NgxSpinnerService,
    useValue: {
      show: jasmine.createSpy('show'),
      hide: jasmine.createSpy('hide'),
      getSpinner: jasmine.createSpy('getSpinner').and.returnValue(EMPTY)
    }
  },
  {
    provide: ToastrService,
    useValue: {
      success: jasmine.createSpy('success'),
      error: jasmine.createSpy('error'),
      warning: jasmine.createSpy('warning'),
      info: jasmine.createSpy('info')
    }
  }
];
