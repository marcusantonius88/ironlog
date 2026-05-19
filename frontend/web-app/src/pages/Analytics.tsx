import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const mockData = [
  { date: 'Week 1', squat: 100, bench: 80, deadlift: 120 },
  { date: 'Week 2', squat: 105, bench: 82, deadlift: 122 },
  { date: 'Week 3', squat: 110, bench: 85, deadlift: 125 },
  { date: 'Week 4', squat: 108, bench: 84, deadlift: 123 },
];

export default function Analytics() {
  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Analytics & Progress</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-gray-800 p-6 rounded border border-gray-700">
          <h3 className="text-amber-400 font-semibold text-lg mb-2">Personal Records</h3>
          <p className="text-3xl font-bold">12</p>
        </div>
        <div className="bg-gray-800 p-6 rounded border border-gray-700">
          <h3 className="text-amber-400 font-semibold text-lg mb-2">Total Volume (Week)</h3>
          <p className="text-3xl font-bold">45.5K lbs</p>
        </div>
        <div className="bg-gray-800 p-6 rounded border border-gray-700">
          <h3 className="text-amber-400 font-semibold text-lg mb-2">Workouts (Month)</h3>
          <p className="text-3xl font-bold">16</p>
        </div>
      </div>

      <div className="bg-gray-800 p-6 rounded border border-gray-700">
        <h2 className="text-xl font-semibold mb-4">Load Progression</h2>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={mockData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#444" />
            <XAxis dataKey="date" stroke="#999" />
            <YAxis stroke="#999" />
            <Tooltip contentStyle={{ backgroundColor: '#1f2937', border: '1px solid #666' }} />
            <Legend />
            <Line type="monotone" dataKey="squat" stroke="#fbbf24" strokeWidth={2} />
            <Line type="monotone" dataKey="bench" stroke="#60a5fa" strokeWidth={2} />
            <Line type="monotone" dataKey="deadlift" stroke="#34d399" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
