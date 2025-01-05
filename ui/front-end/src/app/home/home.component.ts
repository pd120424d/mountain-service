import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router'; 

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
  standalone: true,  // Mark as standalone
  imports: [RouterModule, CommonModule]})
export class HomeComponent {
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

  currentImageIndex: number;

  constructor() {
    this.currentImageIndex = Math.floor(Math.random() * this.images.length);

    this.startSlideshow();
  }

  startSlideshow(): void {
    setInterval(() => {
      this.currentImageIndex = (this.currentImageIndex + 1) % this.images.length;
    }, 5000); // Change image every 5 seconds
  }
}
