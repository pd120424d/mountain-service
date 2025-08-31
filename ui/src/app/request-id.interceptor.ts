import { HttpInterceptorFn, HttpRequest, HttpHandlerFn } from '@angular/common/http';

const HEADER = 'X-Request-ID';

function generateRequestId(): string {
  const array = new Uint8Array(16);
  crypto.getRandomValues(array);
  return Array.from(array).map(b => b.toString(16).padStart(2, '0')).join('');
}

export const requestIdInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn) => {
  // Ensure a unique request ID per request unless one is already set
  let requestId = req.headers.get(HEADER);
  if (!requestId) {
    requestId = generateRequestId();
  }

  const cloned = req.clone({ setHeaders: { [HEADER]: requestId } });
  return next(cloned);
};

