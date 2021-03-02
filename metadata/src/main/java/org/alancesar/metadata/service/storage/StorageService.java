package org.alancesar.metadata.service.storage;

import java.io.InputStream;

public interface StorageService {

    InputStream getChunk(String id, long length) throws Exception;
}
