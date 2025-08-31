import { HttpInterceptorFn, HttpRequest, HttpHandlerFn } from '@angular/common/http';

const HEADER = 'X-Request-ID';

function generateRequestId(): string {
  const array = new Uint8Array(16);
  crypto.getRandomValues(array);
  return Array.from(array).map(b => b.toString(16).padStart(2, '0')).join('');
}

export const requestIdInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn) => {
  // Reuse existing header if present (e.g., server-to-server or retries)
  let requestId = req.headers.get(HEADER) || sessionStorage.getItem(HEADER);
  if (!requestId) {
    requestId = generateRequestId();
    sessionStorage.setItem(HEADER, requestId);
  }

  const cloned = req.clone({ setHeaders: { [HEADER]: requestId } });
  return next(cloned);
};

