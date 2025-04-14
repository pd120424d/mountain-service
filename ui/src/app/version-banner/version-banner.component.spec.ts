import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VersionBannerComponent } from './version-banner.component';

describe('VersionBannerComponent', () => {
  let component: VersionBannerComponent;
  let fixture: ComponentFixture<VersionBannerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VersionBannerComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VersionBannerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
