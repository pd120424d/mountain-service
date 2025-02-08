import {
  TranslateLoader,
  TranslateService,
  TranslateStore,
  MissingTranslationHandler,
  TranslateDefaultParser,
  TranslateCompiler,
  TranslateFakeCompiler,
  TranslateParser,
  USE_DEFAULT_LANG,
  DEFAULT_LANGUAGE,
  ISOLATE_TRANSLATE_SERVICE,
  USE_EXTEND,
} from '@ngx-translate/core';
import { HttpClient } from '@angular/common/http';
import { TranslateHttpLoader } from '@ngx-translate/http-loader';
import { Provider } from '@angular/core';

// Factory to load translations
export function createTranslateLoader(http: HttpClient): TranslateLoader {
  return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}

// Custom Missing Translation Handler
export class CustomMissingTranslationHandler implements MissingTranslationHandler {
  handle(params: { key: string }): string {
    return `Missing: ${params.key}`;
  }
}

// Provide TranslateModule configuration
export function provideTranslate(): Provider[] {
  return [
    { provide: TranslateLoader, useFactory: createTranslateLoader, deps: [HttpClient] },
    { provide: TranslateCompiler, useClass: TranslateFakeCompiler },
    { provide: TranslateParser, useClass: TranslateDefaultParser },
    { provide: MissingTranslationHandler, useClass: CustomMissingTranslationHandler },
    { provide: USE_DEFAULT_LANG, useValue: true }, // Enables default language behavior
    { provide: DEFAULT_LANGUAGE, useValue: 'sr-cyr' }, // Set default language
    { provide: ISOLATE_TRANSLATE_SERVICE, useValue: 'sr-cyr' },
    { provide: USE_EXTEND, useValue: 'sr-cyr' },
    TranslateStore,
    TranslateService,
  ];
}
