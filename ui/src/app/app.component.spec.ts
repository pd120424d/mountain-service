import { TestBed } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { sharedTestingProviders } from './test-utils/shared-test-imports';

describe('AppComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppComponent],
      providers: [...sharedTestingProviders],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });
});
