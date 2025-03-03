import React from 'react';
import StockGraph from './stock-graph';
import EmojiScrollbar from './scroll-bar-emoji'
const InnerListInfo = ({ stock }) => {
  return (
    <div className="inner-info">
      <div className="info-header">
        <h3>Detailed Information</h3>
      </div>
      <div className="graph-and-scroll">

   
        <StockGraph symbol={stock}/>
      
      <EmojiScrollbar titletop = {'Positive sentiment'} titlebottom = {'negative sentiment'}symbol={stock} emojiTop={'😁'} emojiBottom={'😡'}/>
      <EmojiScrollbar  titletop = {'Hot topic'} titlebottom = {'Little coverage'}symbol={stock} emojiTop={'🔥'} emojiBottom={'❄'}/>

      </div>

      <style jsx>{`
      .graph-and-scroll {
          display: flex;
          flex-direction: row;
          align-items: center;
          justify-content: space-between;
        }
        .inner-info {
          padding: 16px;
          margin-top: 16px;
          background: var(--color-platinum);
          border-radius: 8px;
          box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.05);
        }

        .info-header {
          margin-bottom: 16px;
        }

        .info-header h3 {
          color: var(--color-midnight-green);
          margin: 0;
          font-size: 1.1em;
        }

        .stock-graph {
          background: white;
          padding: 16px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
      `}</style>
    </div>
  );
};

export default InnerListInfo;