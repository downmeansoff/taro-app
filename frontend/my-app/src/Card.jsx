import { motion } from 'framer-motion';
import { useState } from 'react';
import PropTypes from 'prop-types';

const Card = ({ card, isSelected, onSelect }) => {
  const [isFlipped, setIsFlipped] = useState(false);

  return (
    <motion.div
      layout
      className="card-container"
      initial={{ scale: 0.9, opacity: 0 }}
      animate={{ 
        scale: isSelected ? 1.1 : 1,
        opacity: 1,
        rotateY: isFlipped ? 180 : 0
      }}
      transition={{ duration: 0.5 }}
      style={{
        position: 'relative',
        width: '100px',
        height: '150px',
        cursor: 'pointer',
        perspective: '1000px'
      }}
      onClick={() => {
        if (!isSelected) {
          setIsFlipped(!isFlipped);
          onSelect(card);
        }
      }}
    >
      {/* Передняя сторона (рубашка) */}
      <motion.div
        className="card-front"
        style={{
          position: 'absolute',
          width: '100%',
          height: '100%',
          background: '#2c3e50',
          borderRadius: '10px',
          backfaceVisibility: 'hidden'
        }}
      />
      
      {/* Задняя сторона (изображение карты) */}
      <motion.div
        className="card-back"
        animate={{ rotateY: isFlipped ? 0 : 180 }}
        style={{
          position: 'absolute',
          width: '100%',
          height: '100%',
          background: '#ecf0f1',
          borderRadius: '10px',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          backfaceVisibility: 'hidden'
        }}
      >
        <img 
          src={`/api/images/${card.imageURL}`} 
          alt={card.name}
          style={{ width: '80%', height: 'auto' }}
        />
        <p style={{ fontSize: '0.8rem', margin: '5px 0' }}>{card.name}</p>
      </motion.div>
    </motion.div>
  );
};

Card.propTypes = {
  card: PropTypes.object.isRequired,
  isSelected: PropTypes.bool,
  onSelect: PropTypes.func.isRequired
};

export default Card;