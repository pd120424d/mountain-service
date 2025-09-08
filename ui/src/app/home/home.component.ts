import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AuthService } from '../services/auth.service';

import { NgxSpinnerService } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';
import { NgxSpinnerModule } from 'ngx-spinner';
import { environment } from '../../environments/environment';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
  standalone: true,
  imports: [RouterModule, CommonModule, TranslateModule, NgxSpinnerModule]
})
export class HomeComponent implements OnInit {
  images: string[] = [
    'assets/slika_1.jpg',
    'assets/slika_2.jpg',
    'assets/slika_3.jpg',
    'assets/slika_4.jpg',
    'assets/slika_5.jpg',
    'assets/slika_6.jpg',
    'assets/slika_7.jpg',
    'assets/slika_8.jpg',
    'assets/slika_9.jpg',
    'assets/slika_10.jpg',
    'assets/slika_11.jpg',
  ];

  currentImageIndex = 0;
  prevImageIndex = 0;
  isStaging = environment.staging;



  constructor(
    private translate: TranslateService,
    public authService: AuthService,
    private spinner: NgxSpinnerService,
    private toastr: ToastrService) {
    this.translate.setDefaultLang('sr-cyr')
  }

  ngOnInit(): void {
    this.preloadImage(this.images[1]); // preload second image immediately


    setInterval(() => {
      const nextIndex = (this.currentImageIndex + 1) % this.images.length;
      const nextImageUrl = this.images[nextIndex];

      this.preloadImage(nextImageUrl).then(() => {
        this.prevImageIndex = this.currentImageIndex;
        this.currentImageIndex = nextIndex;
      });
    }, 8000);
  }



  preloadImage(url: string): Promise<void> {
    return new Promise((resolve) => {
      const img = new Image();
      img.src = url;
      img.onload = () => resolve();
    });
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
  }

  testSpinner() {
    console.log('Spinner clicked');
    this.spinner.show(undefined, {
      type: 'ball-scale-multiple',
      size: 'large',
      bdColor: 'rgba(0, 0, 0, 0.8)',
      color: '#fff'
    });

    setTimeout(() => {
      this.spinner.hide();
      this.toastr.success('Spinner done!');
    }, 2000);
  }

  testToastr() {
    console.log('Toastr clicked');
    this.toastr.info('Toastr test!');
  }
}
