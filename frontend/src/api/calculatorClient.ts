import {
  type Operation,
  type CalculationRequest,
  type CalculationResponse,
  type ApiResult,
  type ApiErrorResponse,
} from '../types/calculator.ts';

/**
 * Shared helper for executing POST operations against backend REST endpoints.
 */
async function postOperation<T>(
  endpoint: string,
  body: CalculationRequest
): Promise<ApiResult<T>> {
  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    const data = await response.json().catch(() => null);

    if (!response.ok) {
      if (data && typeof data === 'object' && 'error' in data) {
        const errorData = data as ApiErrorResponse;
        return {
          ok: false,
          error: errorData.error || {
            code: 'HTTP_ERROR',
            message: `HTTP Error ${response.status}`,
          },
        };
      }
      return {
        ok: false,
        error: {
          code: 'HTTP_ERROR',
          message: `Request failed with status ${response.status}`,
        },
      };
    }

    if (
      !data ||
      typeof data !== 'object' ||
      !('result' in data) ||
      typeof (data as { result: unknown }).result !== 'number'
    ) {
      return {
        ok: false,
        error: {
          code: 'INVALID_RESPONSE',
          message: 'Received invalid response structure from server.',
        },
      };
    }

    return { ok: true, value: data as T };
  } catch (err) {
    return {
      ok: false,
      error: {
        code: 'NETWORK_ERROR',
        message:
          err instanceof Error
            ? err.message
            : 'Unable to connect to calculation service.',
      },
    };
  }
}

/**
 * Dedicated API client exposing strongly-typed calculation functions.
 */
export const calculatorClient = {
  add: (a: number, b: number) =>
    postOperation<CalculationResponse>('/api/add', { a, b }),

  subtract: (a: number, b: number) =>
    postOperation<CalculationResponse>('/api/subtract', { a, b }),

  multiply: (a: number, b: number) =>
    postOperation<CalculationResponse>('/api/multiply', { a, b }),

  divide: (a: number, b: number) =>
    postOperation<CalculationResponse>('/api/divide', { a, b }),

  power: (a: number, b: number) =>
    postOperation<CalculationResponse>('/api/power', { a, b }),

  sqrt: (a: number) =>
    postOperation<CalculationResponse>('/api/sqrt', { a }),

  percentage: (a: number, b: number) =>
    postOperation<CalculationResponse>('/api/percentage', { a, b }),

  /**
   * Dispatches to the correct typed method based on the selected operation.
   * Binary operations require `b`; if it's missing this returns a
   * MISSING_FIELD error instead of silently substituting a default value.
   */
  calculate: (
    op: Operation,
    a: number,
    b?: number
  ): Promise<ApiResult<CalculationResponse>> => {
    const missingB = (): Promise<ApiResult<CalculationResponse>> =>
      Promise.resolve({
        ok: false,
        error: {
          code: 'MISSING_FIELD',
          message: `Operation '${op}' requires a second value.`,
        },
      });

    switch (op) {
      case 'add':
        return b === undefined ? missingB() : calculatorClient.add(a, b);
      case 'subtract':
        return b === undefined ? missingB() : calculatorClient.subtract(a, b);
      case 'multiply':
        return b === undefined ? missingB() : calculatorClient.multiply(a, b);
      case 'divide':
        return b === undefined ? missingB() : calculatorClient.divide(a, b);
      case 'power':
        return b === undefined ? missingB() : calculatorClient.power(a, b);
      case 'sqrt':
        return calculatorClient.sqrt(a);
      case 'percentage':
        return b === undefined ? missingB() : calculatorClient.percentage(a, b);
      default:
        return Promise.resolve({
          ok: false,
          error: {
            code: 'UNSUPPORTED_OPERATION',
            message: `Unsupported operation: ${op}`,
          },
        });
    }
  },
};