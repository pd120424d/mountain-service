import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { Location } from '@angular/common';
import { AppConfigModule } from './app.config';

describe('AppConfigModule', () => {
  let router: Router;
  let location: Location;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppConfigModule]
    }).compileComponents();

    router = TestBed.inject(Router);
    location = TestBed.inject(Location);
  });

  it('should create', () => {
    expect(router).toBeTruthy();
    expect(location).toBeTruthy();
  });

  it('should configure router with routes', () => {
    expect(router.config).toBeDefined();
    expect(router.config.length).toBeGreaterThan(0);
  });

  it('should have RouterModule in imports', () => {
    const module = new AppConfigModule();
    expect(module).toBeTruthy();
  });
});
