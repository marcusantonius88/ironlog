import React, { useState } from 'react';
import DSLParser from '../components/DSLParser';

export default function WorkoutEntry() {
  const [dslText, setDslText] = useState('');
  const [parsedData, setParsedData] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleParse = async () => {
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8081/parse', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ raw_text: dslText }),
      });
      const data = await response.json();
      setParsedData(data);
    } catch (error) {
      console.error('Parsing error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Log Workout</h1>
      
      <div className="grid grid-cols-2 gap-8">
        <div>
          <h2 className="text-xl font-semibold mb-4">DSL Input</h2>
          <textarea
            value={dslText}
            onChange={(e) => setDslText(e.target.value)}
            placeholder="EXERCISE NAME&#10;&#10;Warm up: 1x 1-20 10kg (10)&#10;Work: 3x 6-8 20kg (8 7 6)"
            className="w-full h-64 p-4 bg-gray-800 border border-gray-700 rounded text-white font-mono"
          />
          <button
            onClick={handleParse}
            disabled={loading}
            className="mt-4 w-full bg-amber-600 hover:bg-amber-700 disabled:bg-gray-600 text-white font-bold py-2 px-4 rounded"
          >
            {loading ? 'Parsing...' : 'Parse DSL'}
          </button>
        </div>

        <div>
          <h2 className="text-xl font-semibold mb-4">Preview</h2>
          {parsedData && (
            <DSLParser data={parsedData} />
          )}
        </div>
      </div>
    </div>
  );
}
