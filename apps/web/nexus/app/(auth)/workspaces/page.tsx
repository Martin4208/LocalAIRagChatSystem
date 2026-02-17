'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { workspaceApi } from '@/lib/api-client';
import { Button } from '@/components/ui/button';
import { Header } from '@/components/layout/header';

export interface Workspace {
    id: string;
    name: string;
    description: string | null;
    created_at: string;
    updated_at: string;
}

export default function WorkspacePage() {
    const [workspaces, setWorkspace] = useState<Workspace[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string|null>(null);

    const router = useRouter();

    useEffect(() => {
        const fetchData = async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await workspaceApi.list();
            console.log(data);
            setWorkspace(data.workspaces || []);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load')
        } finally {
            setLoading(false);
        }
        };

        fetchData();
    }, []);

    if (loading) {
        return (
            <div>
            ロード中...
            </div>
        );
        };

        if (error) {
        return (
            <div>
            {error}
            </div>
        );
    };

    const handleNewWorkspace = () => {
        router.push(`/workspaces/new`);
    }

    return (
        <div className="container mx-auto p-8">
            <h1 className="text-3xl font-bold mb-6">Workspaces</h1>

            <Button
                variant="outline"
                onClick={() => handleNewWorkspace()}
            >
                + New Workspace
            </Button>

            {workspaces.length === 0 ? (
                <div>No workspaces</div>
            ) : (
                <div className="grid grid-cols-3 gap-4">
                {workspaces.map(w => (
                    <div key={w.id}>
                        <Button 
                            variant="outline" 
                            onClick={() => router.push(`/workspaces/${w.id}`)}
                        >
                            <h2 className="text-xl font-bold">{w.name}</h2>
                        </Button>
                        <p className="text-gray-600">{w.description}</p>
                        <p className="text-sm text-gray-400">{new Date(w.created_at).toLocaleDateString()}</p>
                    </div>
                ))}
                </div>
            )}
        </div>
    );
}