package org.alancesar.metadata.service.storage;

import java.io.InputStream;

public interface StorageService {

    InputStream get(String id) throws Exception;
}
