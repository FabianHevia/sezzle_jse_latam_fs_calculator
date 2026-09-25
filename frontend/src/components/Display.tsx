import React from 'react';
import { type Operation } from '../types/calculator';

interface DisplayProps {
  operandA: string;
  operandB: string;
  operation: Operation;
  result: number | null;
  errorMessage: string | null;
  isLoading: boolean;
}

const OPERATION_SYMBOLS: Record<Operation, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: '^',
  sqrt: '√',
  percentage: '% of',
};

export const Display: React.FC<DisplayProps> = ({
  operandA,
  operandB,
  operation,
  result,
  errorMessage,
  isLoading,
}) => {
  const formatExpression = (): string => {
    if (!operandA && !operandB && result === null) {
      return 'Ready';
    }

    if (operation === 'sqrt') {
      return `√(${operandA || '0'})`;
    }

    const symbol = OPERATION_SYMBOLS[operation];
    return `${operandA || '0'} ${symbol} ${operandB || '...'}`;
  };

  return (
    <div className="calculator-display" data-testid="calculator-display">
      <div className="expression-preview" aria-live="polite">
        {formatExpression()}
      </div>

      <div className="main-display">
        {isLoading ? (
          <div className="loading-indicator" role="status" aria-label="Calculating">
            <span className="spinner" />
            <span>Calculating...</span>
          </div>
        ) : errorMessage ? (
          <div className="error-message" role="alert">
            {errorMessage}
          </div>
        ) : (
          <div className="result-value" aria-label="Result">
            {result !== null ? result : operandA || '0'}
          </div>
        )}
      </div>
    </div>
  );
};