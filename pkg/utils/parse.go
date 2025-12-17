package utils

import (
	"strings"
)

const (
	defaultDomain    = "docker.io"
	officialRepoName = "library"
)

// ReplaceImageName adds a mirror prefix to container image names, but skips domains in ignoreDomains
func ReplaceImageName(prefix string, ignoreDomains []string, name string) string {
	parts := strings.SplitN(name, "/", 3)
	if parts[0] == prefix {
		return name
	}

	// 统一处理逻辑：检测第一个部分是否是域名
	if len(parts) > 0 && isDomain(parts[0]) {
		// 处理传统域名
		if isLegacyDefaultDomain(parts[0]) {
			parts[0] = defaultDomain
		}
		
		// 检查是否在白名单中
		if shouldIgnoreDomain(parts[0], ignoreDomains) {
			return name // 直接返回原始名称
		}
		
		// 对于域名镜像，统一处理为 prefix/剩余部分
		if len(parts) == 2 {
			// 特殊情况：docker.io/nginx → prefix/library/nginx
			if parts[0] == defaultDomain && parts[1] != officialRepoName {
				return strings.Join([]string{prefix, officialRepoName, parts[1]}, "/")
			}
			return strings.Join([]string{prefix, parts[1]}, "/")
		} else if len(parts) == 3 {
			// docker.io/library/nginx → prefix/library/nginx
			return strings.Join([]string{prefix, parts[1], parts[2]}, "/")
		}
	}
	
	switch len(parts) {
	case 1:
		if shouldIgnoreDomain(defaultDomain, ignoreDomains) {
			return strings.Join([]string{defaultDomain, officialRepoName, parts[0]}, "/")
		}

		return strings.Join([]string{prefix, defaultDomain, officialRepoName, parts[0]}, "/")
	case 2:
		if !isDomain(parts[0]) {
			if shouldIgnoreDomain(defaultDomain, ignoreDomains) {
				return strings.Join([]string{defaultDomain, parts[0], parts[1]}, "/")
			}

			return strings.Join([]string{prefix, defaultDomain, parts[0], parts[1]}, "/")
		}

		if isLegacyDefaultDomain(parts[0]) {
			parts[0] = defaultDomain
		}

		if shouldIgnoreDomain(parts[0], ignoreDomains) {
			return strings.Join([]string{parts[0], parts[1]}, "/")
		}

		return strings.Join([]string{prefix, parts[0], parts[1]}, "/")
	case 3:
		if !isDomain(parts[0]) {
			if shouldIgnoreDomain(defaultDomain, ignoreDomains) {
				return strings.Join([]string{defaultDomain, parts[0], parts[1], parts[2]}, "/")
			}

			return strings.Join([]string{prefix, defaultDomain, parts[0], parts[1], parts[2]}, "/")
		}

		if isLegacyDefaultDomain(parts[0]) {
			parts[0] = defaultDomain
		}

		if shouldIgnoreDomain(parts[0], ignoreDomains) {
			return strings.Join([]string{parts[0], parts[1], parts[2]}, "/")
		}

		return strings.Join([]string{prefix, parts[0], parts[1], parts[2]}, "/")
	}
	return name
}

func isDomain(name string) bool {
	return strings.Contains(name, ".")
}

// shouldIgnoreDomain checks if the image domain should be ignored
func shouldIgnoreDomain(domain string, ignoreDomains []string) bool {
	for _, ignoreDomain := range ignoreDomains {
		if domain == ignoreDomain {
			return true
		}
	}
	return false
}

func isDomain(name string) bool {
	return strings.Contains(name, ".")
}

var (
	legacyDefaultDomain = map[string]struct{}{
		"index.docker.io":      {},
		"registry-1.docker.io": {},
	}
)

func isLegacyDefaultDomain(name string) bool {
	_, ok := legacyDefaultDomain[name]
	return ok
}
