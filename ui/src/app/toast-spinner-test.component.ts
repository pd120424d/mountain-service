import { Component } from '@angular/core';
import { NgxSpinnerService } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';
import { NgIf } from '@angular/common';
import { NgxSpinnerModule } from 'ngx-spinner';

@Component({
  standalone: true,
  selector: 'app-toast-spinner-test',
  template: `
    <div style="padding: 2rem;">
      <h2>Toastr + Spinner Test</h2>
      <button (click)="testSpinner()">Test Spinner</button>
      <button (click)="testToastr()">Test Toastr</button>

      <ngx-spinner></ngx-spinner>
    </div>
  `,
  imports: [NgIf, NgxSpinnerModule]
})
export class ToastSpinnerTestComponent {
  constructor(private spinner: NgxSpinnerService, private toastr: ToastrService) {}

  testSpinner() {
    console.log('Spinner clicked');
    this.spinner.show();
    setTimeout(() => {
      this.spinner.hide();
      this.toastr.success('Spinner finished');
    }, 2000);
  }

  testToastr() {
    console.log('Toastr clicked');
    this.toastr.info('Test toast');
  }
}
