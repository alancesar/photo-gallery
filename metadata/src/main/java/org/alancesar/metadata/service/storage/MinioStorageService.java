package org.alancesar.metadata.service.storage;

import io.minio.GetObjectArgs;
import io.minio.MinioClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.io.InputStream;

@Service
public class MinioStorageService implements StorageService {
    @Value("${minio.bucket.name}")
    private String bucketName;
    private MinioClient client;

    public MinioStorageService(MinioClient client) {
        this.client = client;
    }

    public InputStream get(String filename) throws Exception {
        var args = GetObjectArgs.builder()
                .bucket(bucketName)
                .object(filename)
                .build();

        return client.getObject(args);
    }
}
