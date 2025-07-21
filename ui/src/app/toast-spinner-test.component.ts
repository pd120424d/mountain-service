import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ToastrService } from 'ngx-toastr';
import { NgxSpinnerService } from 'ngx-spinner';
import { NgxSpinnerModule } from 'ngx-spinner';

@Component({
  selector: 'app-toast-spinner-test',
  standalone: true,
  imports: [CommonModule, NgxSpinnerModule],
  template: `
    <div class="test-container">
      <h2>Toast and Spinner Test</h2>
      <div class="button-group">
        <button (click)="testToastr()" class="btn btn-primary">Test Toastr</button>
        <button (click)="testSpinner()" class="btn btn-secondary">Test Spinner</button>
        <button (click)="testBoth()" class="btn btn-success">Test Both</button>
      </div>
    </div>
    <ngx-spinner></ngx-spinner>
  `,
  styles: [`
    .test-container {
      padding: 20px;
      text-align: center;
    }
    .button-group {
      margin-top: 20px;
    }
    .btn {
      margin: 0 10px;
      padding: 10px 20px;
      border: none;
      border-radius: 5px;
      cursor: pointer;
    }
    .btn-primary { background-color: #007bff; color: white; }
    .btn-secondary { background-color: #6c757d; color: white; }
    .btn-success { background-color: #28a745; color: white; }
  `]
})
export class ToastSpinnerTestComponent {
  constructor(
    private toastr: ToastrService,
    private spinner: NgxSpinnerService
  ) {}

  testToastr() {
    console.log('Testing Toastr...');
    this.toastr.success('Success message!', 'Success');
    this.toastr.info('Info message!', 'Info');
    this.toastr.warning('Warning message!', 'Warning');
    this.toastr.error('Error message!', 'Error');
  }

  testSpinner() {
    console.log('Testing Spinner...');
    this.spinner.show(undefined, {
      type: 'ball-scale-multiple',
      size: 'large',
      bdColor: 'rgba(0, 0, 0, 0.8)',
      color: '#fff'
    });
    setTimeout(() => {
      this.spinner.hide();
      this.toastr.success('Spinner test completed!');
    }, 2000);
  }

  testBoth() {
    console.log('Testing both...');
    this.spinner.show(undefined, {
      type: 'ball-scale-multiple',
      size: 'large',
      bdColor: 'rgba(0, 0, 0, 0.8)',
      color: '#fff'
    });
    this.toastr.info('Starting combined test...');

    setTimeout(() => {
      this.spinner.hide();
      this.toastr.success('Combined test completed!');
    }, 3000);
  }
}
