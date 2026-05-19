import React from 'react';
import { Link } from 'react-router-dom';

export default function Navigation() {
  return (
    <nav className="bg-gray-800 border-b border-gray-700">
      <div className="container mx-auto px-4 py-4 flex justify-between items-center">
        <Link to="/" className="text-2xl font-bold text-amber-500">
          IRONLOG
        </Link>
        <div className="space-x-6">
          <Link to="/" className="hover:text-amber-400">
            Dashboard
          </Link>
          <Link to="/workout" className="hover:text-amber-400">
            New Workout
          </Link>
          <Link to="/analytics" className="hover:text-amber-400">
            Analytics
          </Link>
        </div>
      </div>
    </nav>
  );
}
