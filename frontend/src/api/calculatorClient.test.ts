import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { calculatorClient } from './calculatorClient';

describe('calculatorClient', () => {
  const originalFetch = globalThis.fetch;

  beforeEach(() => {
    globalThis.fetch = vi.fn();
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
    vi.restoreAllMocks();
  });

  it('handles successful addition', async () => {
    (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ result: 15 }),
    });

    const res = await calculatorClient.add(10, 5);

    expect(globalThis.fetch).toHaveBeenCalledWith('/api/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ a: 10, b: 5 }),
    });
    expect(res).toEqual({ ok: true, value: { result: 15 } });
  });

  it('handles successful square root operation', async () => {
    (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ result: 9 }),
    });

    const res = await calculatorClient.sqrt(81);

    expect(globalThis.fetch).toHaveBeenCalledWith('/api/sqrt', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ a: 81 }),
    });
    expect(res).toEqual({ ok: true, value: { result: 9 } });
  });

  it('handles API error response (division by zero 400)', async () => {
    (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 400,
      json: async () => ({
        error: {
          code: 'DIVISION_BY_ZERO',
          message: 'cannot divide by zero',
        },
      }),
    });

    const res = await calculatorClient.divide(10, 0);

    expect(res).toEqual({
      ok: false,
      error: {
        code: 'DIVISION_BY_ZERO',
        message: 'cannot divide by zero',
      },
    });
  });

  it('handles non-JSON error response from backend', async () => {
    (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('JSON parse error');
      },
    });

    const res = await calculatorClient.calculate('add', 5, 5);

    expect(res).toEqual({
      ok: false,
      error: {
        code: 'HTTP_ERROR',
        message: 'Request failed with status 500',
      },
    });
  });

  it('handles network connectivity failure', async () => {
    (globalThis.fetch as ReturnType<typeof vi.fn>).mockRejectedValueOnce(
      new Error('Failed to fetch')
    );

    const res = await calculatorClient.multiply(4, 3);

    expect(res).toEqual({
      ok: false,
      error: {
        code: 'NETWORK_ERROR',
        message: 'Failed to fetch',
      },
    });
  });
});