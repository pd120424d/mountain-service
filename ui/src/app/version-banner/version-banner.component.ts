import { Component, OnInit } from '@angular/core';
import { VersionService, VersionInfo } from '../services/version.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-version-banner',
  standalone: true,
  templateUrl: './version-banner.component.html',
  styleUrls: ['./version-banner.component.css'],
  imports: [CommonModule]
})
export class VersionBannerComponent implements OnInit {
  versionInfo?: VersionInfo;

  constructor(private versionService: VersionService) { }

  ngOnInit(): void {
    this.versionService.getVersion().subscribe(info => {
      this.versionInfo = info;
    });
  }
}
