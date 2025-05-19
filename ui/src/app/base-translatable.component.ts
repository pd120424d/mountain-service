import { TranslateService } from '@ngx-translate/core';

export abstract class BaseTranslatableComponent {
  constructor(protected translate: TranslateService) {
    this.translate.setDefaultLang('sr-cyr');
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
  }
}
