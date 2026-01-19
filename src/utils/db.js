const DB_NAME = 'AiChatDB';
const STORE_NAME = 'messages';
const DB_VERSION = 2; // 升级版本以添加索引

export const openDB = () => {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, DB_VERSION);

        request.onupgradeneeded = (event) => {
            const db = event.target.result;
            let store;
            if (!db.objectStoreNames.contains(STORE_NAME)) {
                store = db.createObjectStore(STORE_NAME, { keyPath: 'id' });
            } else {
                store = event.currentTarget.transaction.objectStore(STORE_NAME);
            }
            // 创建用户ID索引，方便根据用户过滤
            if (!store.indexNames.contains('userId')) {
                store.createIndex('userId', 'userId', { unique: false });
            }
        };

        request.onsuccess = (event) => {
            resolve(event.target.result);
        };

        request.onerror = (event) => {
            reject('IndexedDB error: ' + event.target.errorCode);
        };
    });
};

export const saveMessage = async (userId, message) => {
    const db = await openDB();
    return new Promise((resolve, reject) => {
        const transaction = db.transaction([STORE_NAME], 'readwrite');
        const store = transaction.objectStore(STORE_NAME);
        // 保存消息时强制关联用户ID
        const request = store.put({ ...message, userId });

        request.onsuccess = () => resolve();
        request.onerror = () => reject('Save message error');
    });
};

export const getAllMessages = async (userId) => {
    const db = await openDB();
    return new Promise((resolve, reject) => {
        const transaction = db.transaction([STORE_NAME], 'readonly');
        const store = transaction.objectStore(STORE_NAME);
        const index = store.index('userId');
        const request = index.getAll(userId); // 仅获取该用户的消息

        request.onsuccess = (event) => resolve(event.target.result);
        request.onerror = () => reject('Get messages error');
    });
};

export const clearUserMessages = async (userId) => {
    const db = await openDB();
    const messages = await getAllMessages(userId);
    
    return new Promise((resolve, reject) => {
        const transaction = db.transaction([STORE_NAME], 'readwrite');
        const store = transaction.objectStore(STORE_NAME);
        
        // IndexedDB 没有按索引批量删除的简单方法，逐个删除该用户的记录
        let completed = 0;
        if (messages.length === 0) resolve();

        messages.forEach(msg => {
            const request = store.delete(msg.id);
            request.onsuccess = () => {
                completed++;
                if (completed === messages.length) resolve();
            };
            request.onerror = () => reject('Delete message error');
        });
    });
};
