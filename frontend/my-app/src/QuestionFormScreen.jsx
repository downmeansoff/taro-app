import { useState } from 'react';
import { Button, TextField, IconButton, FormHelperText } from '@mui/material';
import ShuffleIcon from '@mui/icons-material/Shuffle';
import axios from 'axios';

export default function QuestionForm() {
    const [question, setQuestion] = useState('');
    const [error, setError] = useState('');
    
    const handleRandomQuestion = () => {
        axios.get('/api/random-question')
            .then(res => setQuestion(res.data.question.slice(0, 200)));
    };

    const handleSubmit = () => {
        if (question.length > 200) {
            setError('Максимум 200 символов');
            return;
        }
        // Отправка вопроса
    };

    return (
        <div style={{ padding: 20 }}>
            <TextField
                value={question}
                onChange={(e) => {
                    if (e.target.value.length <= 200) {
                        setQuestion(e.target.value);
                        setError('');
                    }
                }}
                fullWidth
                placeholder="Введите ваш вопрос"
                error={!!error}
                inputProps={{ maxLength: 200 }}
                InputProps={{
                    endAdornment: (
                        <IconButton onClick={handleRandomQuestion}>
                            <ShuffleIcon />
                        </IconButton>
                    ),
                }}
            />
            <FormHelperText>
                {question.length}/200 символов
            </FormHelperText>
            
            <Button 
                variant="contained" 
                onClick={handleSubmit}
                style={{ marginTop: 15 }}
            >
                Подтвердить
            </Button>
        </div>
    );
}