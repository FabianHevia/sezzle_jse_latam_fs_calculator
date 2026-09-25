import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom/vitest';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Calculator } from './Calculator';
import { calculatorClient } from '../api/calculatorClient';

vi.mock('../api/calculatorClient', () => ({
  calculatorClient: {
    calculate: vi.fn(),
  },
}));

describe('Calculator Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders calculator elements correctly', () => {
    render(<Calculator />);

    expect(screen.getByRole('heading', { name: /REST Calculator/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/First Number \(a\)/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Second Number \(b\)/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Clear/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Calculate result/i })).toBeInTheDocument();
  });

  it('performs successful addition workflow', async () => {
    const user = userEvent.setup();
    (calculatorClient.calculate as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      value: { result: 42 },
    });

    render(<Calculator />);

    const inputA = screen.getByLabelText(/First Number \(a\)/i);
    const inputB = screen.getByLabelText(/Second Number \(b\)/i);
    const calculateBtn = screen.getByRole('button', { name: /Calculate result/i });

    await user.type(inputA, '20');
    await user.type(inputB, '22');
    await user.click(calculateBtn);

    expect(calculatorClient.calculate).toHaveBeenCalledWith('add', 20, 22);

    await waitFor(() => {
      expect(screen.getByLabelText('Result')).toHaveTextContent('42');
    });
  });

  it('switches operations and hides second input for unary sqrt operation', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    const sqrtBtn = screen.getByRole('button', { name: /Square Root/i });
    await user.click(sqrtBtn);

    expect(screen.getByLabelText(/Number \(x\)/i)).toBeInTheDocument();
    expect(screen.queryByLabelText(/Second Number \(b\)/i)).not.toBeInTheDocument();
  });

  it('displays client-side validation error for empty input', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    const calculateBtn = screen.getByRole('button', { name: /Calculate result/i });
    await user.click(calculateBtn);

    expect(screen.getByRole('alert')).toHaveTextContent('Please enter the first number.');
    expect(calculatorClient.calculate).not.toHaveBeenCalled();
  });

  it('displays validation error for division by zero', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    const divideBtn = screen.getByRole('button', { name: /Divide/i });
    await user.click(divideBtn);

    const inputA = screen.getByLabelText(/First Number \(a\)/i);
    const inputB = screen.getByLabelText(/Second Number \(b\)/i);
    const calculateBtn = screen.getByRole('button', { name: /Calculate result/i });

    await user.type(inputA, '10');
    await user.type(inputB, '0');
    await user.click(calculateBtn);

    expect(screen.getByRole('alert')).toHaveTextContent('Cannot divide by zero.');
    expect(calculatorClient.calculate).not.toHaveBeenCalled();
  });

  it('displays API error message returned from backend', async () => {
    const user = userEvent.setup();
    (calculatorClient.calculate as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      error: {
        code: 'DIVISION_BY_ZERO',
        message: 'cannot divide by zero',
      },
    });

    render(<Calculator />);

    const inputA = screen.getByLabelText(/First Number \(a\)/i);
    const inputB = screen.getByLabelText(/Second Number \(b\)/i);

    await user.type(inputA, '10');
    await user.type(inputB, '2');

    const multiplyBtn = screen.getByRole('button', { name: /Multiply/i });
    await user.click(multiplyBtn);

    const calculateBtn = screen.getByRole('button', { name: /Calculate result/i });
    await user.click(calculateBtn);

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('cannot divide by zero');
    });
  });

  it('shows loading state while awaiting API response', async () => {
    const user = userEvent.setup();
    let resolveApi: (val: unknown) => void;
    const pendingPromise = new Promise((resolve) => {
      resolveApi = resolve;
    });

    (calculatorClient.calculate as ReturnType<typeof vi.fn>).mockReturnValueOnce(pendingPromise);

    render(<Calculator />);

    await user.type(screen.getByLabelText(/First Number \(a\)/i), '5');
    await user.type(screen.getByLabelText(/Second Number \(b\)/i), '5');
    await user.click(screen.getByRole('button', { name: /Calculate result/i }));

    expect(screen.getByRole('status', { name: /Calculating/i })).toBeInTheDocument();

    resolveApi!({ ok: true, value: { result: 10 } });

    await waitFor(() => {
      expect(screen.queryByRole('status')).not.toBeInTheDocument();
      expect(screen.getByLabelText('Result')).toHaveTextContent('10');
    });
  });

  it('clears input fields when Clear button is clicked', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    const inputA = screen.getByLabelText(/First Number \(a\)/i) as HTMLInputElement;
    await user.type(inputA, '123');
    expect(inputA.value).toBe('123');

    await user.click(screen.getByRole('button', { name: /Clear/i }));
    expect(inputA.value).toBe('');
  });
});