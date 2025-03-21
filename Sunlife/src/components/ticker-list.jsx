import React, { useState } from 'react'
import { BsFillChatRightTextFill } from "react-icons/bs";
import TickerGraph from './ticker-graph';
import { useSymbol } from '../contexts/symbol-context';
import watchlistStocks from '../data/watchlist'
const stockData = watchlistStocks
import InnerListInfo from './inner-list-info'

const TickerList = () => {
  const [expandedStock, setExpandedStock] = useState(null);
  const {selectedSymbol, setSelectedSymbol} = useSymbol();

  const handleClickItem = (symbol) => {
    setExpandedStock(expandedStock === symbol ? null : symbol);
  }

  const handleClickChat = (event, symbol) => {
    event.stopPropagation();
    setSelectedSymbol(symbol);
  }

  return (
    <div className="ticker-container">
      {stockData.map((stock) => (
        <div key={stock.symbol} className="stock-item" onClick={() => handleClickItem(stock.symbol)}>
          <div className="stock-details">
            <div className="stock-main-info">
              <div className="stock-identifier">
                <span className="stock-symbol"> {stock.symbol}{stock.emoji}</span>
                <span className="stock-name">{stock.name}</span>
                <span className="stock-sector">{stock.sector}</span>
              </div>
              <div className="stock-price-details">
                <span className="stock-price">${stock.price}</span>
                <span className={`stock-change ${stock.change.startsWith('+') ? 'positive' : 'negative'}`}>
                  {stock.change}
                </span>

                <BsFillChatRightTextFill 
                  className="chat-icon" 
                  onClick={(event) => handleClickChat(event, stock.symbol)}
                />
              </div>
            </div>
          </div>
          {expandedStock === stock.symbol && (
              <InnerListInfo stock ={stock.symbol}/>
          )}
        </div>
      ))}
      <style jsx>{`
        .ticker-container {
          display: flex;
          flex-direction: column;
          gap: 12px;
          padding: 16px;
        }

        .stock-item {
          background: white;
          padding: 16px;
          border-radius: 12px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
          transition: all 0.3s ease;
          cursor: pointer;
          width: 100%;
        }

        .stock-details {
          width: 100%;
        }

        .stock-main-info {
          display: flex;
          justify-content: space-between;
          align-items: center;
          width: 100%;
        }

        .stock-identifier {
          display: flex;
          align-items: center;
          gap: 16px;
          flex: 2;
        }

        .stock-symbol {
          font-size: 1.2em;
          font-weight: bold;
          color: var(--color-midnight-green);
          min-width: 120px;
        }

        .stock-name {
          color: #666;
          font-size: 0.9em;
          min-width: 150px;
        }

        .stock-sector {
          padding: 4px 8px;
          background: var(--color-platinum);
          border-radius: 4px;
          font-size: 0.8em;
        }

        .stock-price-details {
          display: flex;
          align-items: center;
          gap: 16px;
          flex: 1;
          justify-content: flex-end;
        }

        .stock-price {
          font-weight: bold;
          color: var(--color-midnight-green);
        }

        .stock-change.positive {
          color: #22c55e;
        }

        .stock-change.negative {
          color: #ef4444;
        }

        .stock-trend {
          font-size: 0.8em;
          padding: 2px 6px;
          border-radius: 4px;
        }

        .stock-trend.bullish {
          background: #dcfce7;
          color: #15803d;
        }

        .stock-trend.bearish {
          background: #fee2e2;
          color: #b91c1c;
        }

        .stock-trend.neutral {
          background: #f3f4f6;
          color: #4b5563;
        }

        .chat-icon {
          cursor: pointer;
          color: var(--color-midnight-green);
          transition: transform 0.2s ease;
        }

        .chat-icon:hover {
          transform: scale(1.2);
        }

        .stock-graph {
          margin-top: 16px;
          padding: 0;
          background: var(--color-platinum);
          border-radius: 8px;
        }
      `}</style>
    </div>
  );
};

export default TickerList;