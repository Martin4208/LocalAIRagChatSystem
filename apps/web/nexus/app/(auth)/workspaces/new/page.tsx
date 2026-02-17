// apps/web/nexus/app/(auth)/workspaces/new/page.tsx
'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { useCreateWorkspace } from '@/lib/hooks/use-workspaces';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';

export default function NewWorkspace() {
    const [name, setName] = useState('');
    const [ description, setDescription ] = useState('');
    
    const router = useRouter();
    const createMutation = useCreateWorkspace();

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        if (!name.trim()) {
            return;
        }

        try {
            const result = await createMutation.mutateAsync({
                name,
                description: description || undefined,
            });
            const workspaceId = result.id;
            router.push(`/workspaces/${workspaceId}`);
        } catch (error) {

        }
    }

    return (
        <div>
            <Card>
                <CardHeader>
                    <CardTitle>Create New Workspace</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit}>
                        <div>
                            <label>Workspace Name *</label>
                            <Input 
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                            />
                            <div>
                                <label>Description (optional)</label>
                                <Textarea 
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                />
                            </div>
                            <Button
                                type="submit"
                                disabled={!name.trim() || createMutation.isPending}
                            >
                                {createMutation.isPending ? 'Creating...' : 'New Workspace'}
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
}