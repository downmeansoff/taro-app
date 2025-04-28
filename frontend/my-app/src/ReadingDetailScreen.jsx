import { useParams } from 'react-router-dom';
import { useEffect, useState } from 'react';
import axios from 'axios';

export default function ReadingDetail() {
    const { id } = useParams();
    const [reading, setReading] = useState(null);

    useEffect(() => {
        axios.get(`/api/history/${id}`)
            .then(res => setReading(res.data))
            .catch(err => console.error(err));
    }, [id]);

    if (!reading) return <div>Loading...</div>;

    return (
        <div>
            <h3>{reading.question}</h3>
            <p>{reading.answer}</p>
            <p>Карты: {reading.cards.join(', ')}</p>
        </div>
    );
}