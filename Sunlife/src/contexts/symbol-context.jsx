import React, { createContext, useContext, useState } from 'react';

const SymbolContext = createContext();

const SymbolProvider = ({ children }) => {
    const [selectedSymbol, setSelectedSymbol] = useState('');

    return (
        <SymbolContext.Provider value={{ selectedSymbol, setSelectedSymbol }}>
            {children}
        </SymbolContext.Provider>
    );
};

const useSymbol = () => {
    return useContext(SymbolContext)
}

export { SymbolProvider, useSymbol };