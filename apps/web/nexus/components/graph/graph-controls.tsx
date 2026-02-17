'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Search, Box, Layers, X } from 'lucide-react';
import type { GraphFilter } from '@/lib/types/graph';

interface GraphControlsProps {
  mode: '2d' | '3d';
  onModeChange: (mode: '2d' | '3d') => void;
  filter: GraphFilter;
  onFilterChange: (filter: GraphFilter) => void;
  availableNodeTypes: string[];
  availableEdgeTypes: string[];
}

export function GraphControls({
  mode,
  onModeChange,
  filter,
  onFilterChange,
  availableNodeTypes,
  availableEdgeTypes,
}: GraphControlsProps) {
  const handleSearchChange = (value: string) => {
    onFilterChange({ ...filter, searchQuery: value });
  };

  const handleNodeTypeToggle = (type: string) => {
    const current = filter.nodeTypes || [];
    const updated = current.includes(type)
      ? current.filter((t) => t !== type)
      : [...current, type];
    onFilterChange({ ...filter, nodeTypes: updated });
  };

  const handleClearFilters = () => {
    onFilterChange({
      nodeTypes: [],
      edgeTypes: [],
      searchQuery: '',
    });
  };

  const hasActiveFilters =
    (filter.nodeTypes && filter.nodeTypes.length > 0) ||
    (filter.edgeTypes && filter.edgeTypes.length > 0) ||
    (filter.searchQuery && filter.searchQuery.trim());

  return (
    <div className="absolute top-4 left-4 z-10 bg-white rounded-lg shadow-lg border p-4 w-80 max-h-[calc(100vh-8rem)] overflow-y-auto">
      <div className="space-y-4">
        {/* Mode Toggle */}
        <div>
          <Label className="text-sm font-medium mb-2 block">View Mode</Label>
          <div className="flex gap-2">
            <Button
              variant={mode === '2d' ? 'default' : 'outline'}
              size="sm"
              onClick={() => onModeChange('2d')}
              className="flex-1"
            >
              <Layers className="h-4 w-4 mr-2" />
              2D
            </Button>
            <Button
              variant={mode === '3d' ? 'default' : 'outline'}
              size="sm"
              onClick={() => onModeChange('3d')}
              className="flex-1"
            >
              <Box className="h-4 w-4 mr-2" />
              3D
            </Button>
          </div>
        </div>

        {/* Search */}
        <div>
          <Label htmlFor="graph-search" className="text-sm font-medium mb-2 block">
            Search Nodes
          </Label>
          <div className="relative">
            <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              id="graph-search"
              placeholder="Search by label..."
              value={filter.searchQuery || ''}
              onChange={(e) => handleSearchChange(e.target.value)}
              className="pl-8"
            />
          </div>
        </div>

        {/* Node Type Filter */}
        <div>
          <div className="flex items-center justify-between mb-2">
            <Label className="text-sm font-medium">Node Types</Label>
            {hasActiveFilters && (
              <Button
                variant="ghost"
                size="sm"
                onClick={handleClearFilters}
                className="h-6 px-2 text-xs"
              >
                <X className="h-3 w-3 mr-1" />
                Clear
              </Button>
            )}
          </div>
          <div className="flex flex-wrap gap-2">
            {availableNodeTypes.length === 0 ? (
              <p className="text-sm text-muted-foreground">No nodes yet</p>
            ) : (
              availableNodeTypes.map((type) => {
                const isActive = filter.nodeTypes?.includes(type) ?? false;
                return (
                  <Badge
                    key={type}
                    variant={isActive ? 'default' : 'outline'}
                    className="cursor-pointer"
                    onClick={() => handleNodeTypeToggle(type)}
                  >
                    {type}
                  </Badge>
                );
              })
            )}
          </div>
        </div>

        {/* Stats */}
        <div className="pt-2 border-t">
          <div className="text-xs text-muted-foreground space-y-1">
            <div className="flex justify-between">
              <span>Total Node Types:</span>
              <span className="font-medium">{availableNodeTypes.length}</span>
            </div>
            <div className="flex justify-between">
              <span>Total Edge Types:</span>
              <span className="font-medium">{availableEdgeTypes.length}</span>
            </div>
            {filter.nodeTypes && filter.nodeTypes.length > 0 && (
              <div className="flex justify-between text-blue-600">
                <span>Filtered Types:</span>
                <span className="font-medium">{filter.nodeTypes.length}</span>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}