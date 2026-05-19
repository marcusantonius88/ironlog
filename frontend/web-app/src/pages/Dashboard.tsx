import React from 'react';

export default function Dashboard() {
  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-4xl font-bold mb-8">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-gradient-to-br from-amber-600 to-amber-700 p-6 rounded-lg shadow-lg">
          <h3 className="text-amber-100 text-sm font-semibold mb-2">Current Streak</h3>
          <p className="text-4xl font-bold">8 Days</p>
        </div>
        <div className="bg-gradient-to-br from-blue-600 to-blue-700 p-6 rounded-lg shadow-lg">
          <h3 className="text-blue-100 text-sm font-semibold mb-2">Personal Records</h3>
          <p className="text-4xl font-bold">12</p>
        </div>
        <div className="bg-gradient-to-br from-green-600 to-green-700 p-6 rounded-lg shadow-lg">
          <h3 className="text-green-100 text-sm font-semibold mb-2">This Month</h3>
          <p className="text-4xl font-bold">16</p>
          <p className="text-green-100 text-xs">Workouts</p>
        </div>
        <div className="bg-gradient-to-br from-purple-600 to-purple-700 p-6 rounded-lg shadow-lg">
          <h3 className="text-purple-100 text-sm font-semibold mb-2">Volume (Week)</h3>
          <p className="text-4xl font-bold">45.5K</p>
          <p className="text-purple-100 text-xs">lbs</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700">
          <h2 className="text-xl font-semibold mb-4">Recent Workouts</h2>
          <div className="space-y-4">
            {[1, 2, 3].map(i => (
              <div key={i} className="bg-gray-700 p-4 rounded flex justify-between items-center">
                <div>
                  <p className="font-semibold">Upper Body A</p>
                  <p className="text-sm text-gray-400">45 min • 4 exercises</p>
                </div>
                <span className="text-amber-400 font-bold">+2,350 lbs</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700">
          <h2 className="text-xl font-semibold mb-4">Recommendations</h2>
          <div className="space-y-4">
            <div className="bg-green-900 bg-opacity-30 p-4 rounded border border-green-700">
              <p className="font-semibold text-green-400">Bench Press</p>
              <p className="text-sm text-gray-300">Try increasing weight to 185 lbs</p>
            </div>
            <div className="bg-yellow-900 bg-opacity-30 p-4 rounded border border-yellow-700">
              <p className="font-semibold text-yellow-400">Squat</p>
              <p className="text-sm text-gray-300">Consider deload week</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
