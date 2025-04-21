import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VersionBannerComponent } from './version-banner.component';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';

describe('VersionBannerComponent', () => {
  let component: VersionBannerComponent;
  let fixture: ComponentFixture<VersionBannerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VersionBannerComponent],
      providers: [...sharedTestingProviders]
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
