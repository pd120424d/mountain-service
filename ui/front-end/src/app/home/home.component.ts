// src/app/home/home.component.ts
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router'; 

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
  standalone: true,  // Mark as standalone
  imports: [RouterModule]
})
export class HomeComponent {
    
  testClick(): void {
    console.log('Button clicked!');
  }
 }
