/**
 * Supported calculator operations.
 * Includes base arithmetic and optional operations (power, sqrt, percentage).
 */
export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage';

/**
 * Standard operation request body.
 * `a` is required for all operations. `b` is required for binary operations.
 */
export interface CalculationRequest {
  a: number;
  b?: number;
}

/**
 * Successful response payload from backend calculation API.
 */
export interface CalculationResponse {
  result: number;
}

/**
 * Machine-readable and human-readable API error structure.
 */
export interface ApiError {
  code: string;
  message: string;
}

/**
 * Error envelope payload from backend API.
 */
export interface ApiErrorResponse {
  error: ApiError;
}

/**
 * Discriminated union representing either a successful API call or an error.
 */
export type ApiResult<T> =
  | { ok: true; value: T }
  | { ok: false; error: ApiError };