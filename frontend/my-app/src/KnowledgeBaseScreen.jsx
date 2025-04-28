import { useState, useEffect } from 'react';
import { TextField, List, ListItem, ListItemText } from '@mui/material';
import axios from 'axios';
import CardModal from './CardModal';

export default function KnowledgeBase() {
  const [cards, setCards] = useState([]);
  const [search, setSearch] = useState('');
  const [selectedCard, setSelectedCard] = useState(null);
  const [modalOpen, setModalOpen] = useState(false);

  useEffect(() => {
    axios.get('/api/cards')
      .then(res => setCards(res.data))
      .catch(console.error);
  }, []);

  const filteredCards = cards.filter(card =>
    card.name.toLowerCase().includes(search.toLowerCase()) ||
    card.keywords.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div style={{ padding: 20 }}>
      <TextField
        label="Поиск по картам"
        fullWidth
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        style={{ marginBottom: 20 }}
      />
      
      <List>
        {filteredCards.map(card => (
          <ListItem 
            button 
            key={card.id}
            onClick={() => {
              setSelectedCard(card);
              setModalOpen(true);
            }}
          >
            <ListItemText 
              primary={card.name} 
              secondary={card.keywords} 
            />
          </ListItem>
        ))}
      </List>

      <CardModal 
        card={selectedCard} 
        open={modalOpen} 
        onClose={() => setModalOpen(false)}
      />
    </div>
  );
}