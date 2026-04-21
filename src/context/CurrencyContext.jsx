import React, { createContext, useState, useContext } from 'react';

const CurrencyContext = createContext();

export const currencies = {
  USD: {
    symbol: '$',
    code: 'USD',
    rate: 1, // Base currency
    format: (val) => `$${val.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
  },
  IDR: {
    symbol: 'Rp',
    code: 'IDR',
    rate: 15000, // Fixed rate for simulation
    format: (val) => `Rp${(val).toLocaleString('id-ID', { minimumFractionDigits: 0, maximumFractionDigits: 0 })}`
  }
};

export function CurrencyProvider({ children }) {
  const [currency, setCurrency] = useState('IDR'); // Default to IDR as requested

  const currentCurrency = currencies[currency];

  const toggleCurrency = () => {
    setCurrency(prev => prev === 'USD' ? 'IDR' : 'USD');
  };

  const formatValue = (usdValue) => {
    const converted = usdValue * currentCurrency.rate;
    return currentCurrency.format(converted);
  };

  return (
    <CurrencyContext.Provider value={{ currency, currentCurrency, toggleCurrency, formatValue }}>
      {children}
    </CurrencyContext.Provider>
  );
}

export function useCurrency() {
  return useContext(CurrencyContext);
}
