'use client';

import { useState, useRef, useEffect } from 'react';

import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export default Search({ documents }) {
    

    return (
        <Card>
            <Input 

            />
            <Button
                type="submit"
                disabled={}
            >
                Search
            </Button>
            <div>
                <p>Filter by concept</p>
                
            </div>
            <div>
                <p>Filter by file type</p>
            </div>
        </Card>
    );
}