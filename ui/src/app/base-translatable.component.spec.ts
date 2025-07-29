import { TestBed } from '@angular/core/testing';
import { TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from './base-translatable.component';
import { Component } from '@angular/core';

// Create a concrete implementation for testing
@Component({
  template: '',
  standalone: true
})
class TestTranslatableComponent extends BaseTranslatableComponent {
  constructor(translate: TranslateService) {
    super(translate);
  }
}

describe('BaseTranslatableComponent', () => {
  let component: TestTranslatableComponent;
  let translateService: jasmine.SpyObj<TranslateService>;

  beforeEach(() => {
    const translateSpy = jasmine.createSpyObj('TranslateService', ['setDefaultLang', 'use']);

    TestBed.configureTestingModule({
      imports: [TestTranslatableComponent],
      providers: [
        { provide: TranslateService, useValue: translateSpy }
      ]
    });

    const fixture = TestBed.createComponent(TestTranslatableComponent);
    component = fixture.componentInstance;
    translateService = TestBed.inject(TranslateService) as jasmine.SpyObj<TranslateService>;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set default language to sr-cyr on construction', () => {
    expect(translateService.setDefaultLang).toHaveBeenCalledWith('sr-cyr');
  });

  it('should switch language when switchLanguage is called', () => {
    const testLanguage = 'en';
    component.switchLanguage(testLanguage);
    expect(translateService.use).toHaveBeenCalledWith(testLanguage);
  });
});
