import { useState } from 'react';
import { Button } from '@mui/material';
import { motion, AnimatePresence } from 'framer-motion';
import Card from '../components/Card';
import axios from 'axios';

const DailyReadingScreen = () => {
  const [cards, setCards] = useState([]);
  const [selectedCard, setSelectedCard] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  // Загрузка карт при монтировании
  useState(() => {
    axios.get('/api/cards')
      .then(res => {
        setCards(res.data.map(c => ({ ...c, id: Math.random().toString() })));
        setIsLoading(false);
      });
  }, []);

  const shuffleCards = () => {
    setCards([...cards.sort(() => Math.random() - 0.5)]);
  };

  const handleCardSelect = (card) => {
    if (!selectedCard) {
      setSelectedCard(card);
    }
  };

  return (
    <div style={{ padding: '20px' }}>
      <Button 
        variant="contained" 
        onClick={shuffleCards}
        style={{ marginBottom: '20px' }}
      >
        Перемешать карты
      </Button>

      {isLoading ? (
        <div>Загрузка карт...</div>
      ) : (
        <motion.div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(100px, 1fr))',
            gap: '15px',
            padding: '10px'
          }}
        >
          <AnimatePresence>
            {cards.map((card) => (
              <Card
                key={card.id}
                card={card}
                isSelected={selectedCard?.id === card.id}
                onSelect={handleCardSelect}
              />
            ))}
          </AnimatePresence>
        </motion.div>
      )}
    </div>
  );
};

export default DailyReadingScreen;