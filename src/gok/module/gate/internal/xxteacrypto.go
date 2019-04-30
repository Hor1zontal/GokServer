/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/12/20
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package internal

import "aliens/common/cipher/xxtea"

func NewXXTeaCrypto(key string) *XXTeaCrypto {
	return &XXTeaCrypto{
		key: []byte(key),
	}
}

type XXTeaCrypto struct {
	key []byte
}

//加密方法
func (this *XXTeaCrypto) Encrypt(data []byte) []byte {
	return xxtea.Encrypt(data, this.key)
}

//解密方法
func (this *XXTeaCrypto) Decrypt(data []byte) []byte {
	return xxtea.Decrypt(data, this.key)
}
