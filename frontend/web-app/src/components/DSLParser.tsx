import React from 'react';

interface SetGroup {
  set_type: string;
  planned_sets: number;
  target_rep_min: number;
  target_rep_max: number;
  weight: number;
  unit: string;
  executed_reps?: number[];
}

interface DSLParserProps {
  data: {
    success: boolean;
    exercise?: string;
    set_groups?: SetGroup[];
    error?: string;
  };
}

export default function DSLParser({ data }: DSLParserProps) {
  if (!data.success) {
    return (
      <div className="bg-red-900 border border-red-700 p-4 rounded text-red-100">
        Error: {data.error}
      </div>
    );
  }

  return (
    <div className="bg-gray-800 border border-gray-700 rounded p-6 space-y-4">
      <div>
        <h3 className="text-lg font-semibold text-amber-400">{data.exercise}</h3>
      </div>

      {data.set_groups?.map((group, idx) => (
        <div key={idx} className="bg-gray-700 p-4 rounded">
          <div className="flex justify-between items-center mb-2">
            <span className="font-semibold text-amber-300">{group.set_type}</span>
            <span className="text-sm text-gray-300">{group.planned_sets}x {group.target_rep_min}-{group.target_rep_max}</span>
          </div>
          <div className="text-sm text-gray-400">
            {group.weight}{group.unit}
          </div>
          {group.executed_reps && (
            <div className="text-sm text-green-400 mt-2">
              Executed: {group.executed_reps.join(', ')} reps
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
