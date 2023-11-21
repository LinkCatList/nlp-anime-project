import pandas as pd
import numpy as np
import nltk
import string
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
import tensorflow as tf

nltk.download('punkt')

train_df = pd.read_csv('/content/anime-description.csv')
import random
train_df['genres'] = train_df['genres'].astype(str)

train_df['genres'] = train_df['genres'].str.split(', ')
train_df['first_genre'] = train_df['genres'].apply(lambda x: x[0])
train_df['second_genre'] = train_df['genres'].apply(lambda x: x[1] if len(x) > 1 else x[0])
train_df = train_df.drop(columns=['genres'])

train_df['Id'] = pd.Series([i for i in range(train_df.shape[0])])
train_df

genre_count_df = train_df[['Id','first_genre']].groupby('first_genre').agg('count').sort_values('Id', ascending = False)

genre_count_df.reset_index(inplace = True)
genre_count_df.rename(columns ={'Id':'genre_cnt'}, inplace = True)

genre_count_df.plot.bar(y = 'genre_cnt', x = 'first_genre')

genre_count_df = train_df[['Id','second_genre']].groupby('second_genre').agg('count').sort_values('Id', ascending = False)

genre_count_df.reset_index(inplace = True)
genre_count_df.rename(columns ={'Id':'genre_cnt'}, inplace = True)

genre_count_df.plot.bar(y = 'genre_cnt', x = 'second_genre')

genre_count_df

train_df = train_df.merge(genre_count_df, how = 'left', left_on='second_genre', right_on='second_genre')

train_df

min_cnt = 700
train_df = train_df[train_df['genre_cnt'] > min_cnt]
train_df['genre_cnt'] = train_df['genre_cnt'].replace('nan', 0)

train_df

test_df =  pd.read_csv('/content/anime-description.csv')# тут нужно импортнуть тестовые данные
test_df['genres'] = test_df['genres'].astype(str)

test_df['genres'] = test_df['genres'].str.split(', ')
test_df['first_genre'] = test_df['genres'].apply(lambda x: x[0])
test_df['second_genre'] = test_df['genres'].apply(lambda x: x[1] if len(x) > 1 else x[0])
test_df = test_df.drop(columns=['genres'])

test_df['Id'] = pd.Series([i for i in range(test_df.shape[0])])

test_df = test_df.merge(genre_count_df, how = 'left', left_on='second_genre', right_on='second_genre')

test_df

def tokenize_text(text):
  tokenized_text = nltk.word_tokenize(text)
  tokens = [i.lower() for i in tokenized_text if ( i not in string.punctuation )]
  return tokens

train_df['description_tokenized'] = train_df['description'].apply(lambda x:tokenize_text(x))
test_df['description_tokenized'] = test_df['description'].apply(lambda x:tokenize_text(x))

test_df

vocab = {}
mx = 1000000

vocab["<PAD>"] = mx+2;
vocab["<START>"] = mx+1;
vocab["<UNKW>"] = mx;

for tokens in train_df['description_tokenized']:
  for word in tokens:
    if word not in vocab.keys():
      vocab[word] = 1
    else:
      vocab[word]+=1
vocab = {k: v for k, v in sorted(vocab.items(), key = lambda item: item[1], reverse = True)}

cnt = 0
for k in vocab.keys():
  vocab[k] = cnt
  cnt+=1
print(len(vocab))
print()
vocab

mini_vocab = {}
for k, v in vocab.items():
  if v<10000:
    mini_vocab[k] = v;
mini_vocab

def EncodeTokens(ListTokens, VocabTokens):
  res = []
  res = [VocabTokens.get(word, VocabTokens['<UNKW>']) for word in ListTokens]
  return [VocabTokens['<START>']]+res
def DecodeTokens(EncodedTokens, VocabTokens):
  res = []
  for i in EncodedTokens:
    for word, ind in VocabTokens.items():
      if i == ind:
        res.append(word)
        break
  return res

train_df['description_encoded'] = train_df['description_tokenized'].apply(lambda x: EncodeTokens(x, mini_vocab))
train_df

test_df['description_encoded'] = test_df['description_tokenized'].apply(lambda x: EncodeTokens(x, mini_vocab))

test_df

train_data = train_df.description_encoded.to_numpy()
train_first_label = pd.get_dummies(train_df['first_genre']).values
train_second_label = pd.get_dummies(train_df['second_genre']).values

test_data = test_df.description_encoded.to_numpy()
test_first_label = pd.get_dummies(test_df['first_genre']).values
test_second_label = pd.get_dummies(test_df['second_genre']).values

train_df['description_len'] = train_df['description_encoded'].apply (len)

print ('минимальная длина описания:', train_df.description_len.min())
print ('средняя длина описания:', round(train_df.description_len.mean()))
print ('максимальная длина описания:', train_df.description_len.max())

plt.hist(train_df.description_len, density = True)

MAX_SEQ_LEN = 70

train_data = tf.keras.preprocessing.sequence.pad_sequences(
    train_data,
    value= vocab['<PAD>'],
    padding= 'post',
    maxlen= MAX_SEQ_LEN)

test_data = tf.keras.preprocessing.sequence.pad_sequences(
    test_data,
    value= vocab['<PAD>'],
    padding= 'post',
    maxlen= MAX_SEQ_LEN)
print('Тернировочные данные:')
print(train_data.shape)
print(train_data[0])
print()
print('Тестовые данные:')
print(test_data.shape)
print(test_data[0])

partial_x_train, x_val, partial_y_train, y_val = train_test_split(train_data, train_second_label,
                                                                  test_size = 0.05, random_state = 42)


print(partial_x_train.shape, partial_y_train.shape)
print(x_val.shape, y_val.shape)

VOCAB_SIZE = len(mini_vocab)
EMB_SIZE = 32
CLASS_NUM = y_val.shape[1]

model = tf.keras.Sequential([
    tf.keras.layers.Embedding(VOCAB_SIZE, EMB_SIZE),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=True, dropout=0.3, recurrent_dropout=0.5)),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=False, dropout=0.3, recurrent_dropout=0.5)),
    tf.keras.layers.Dense(CLASS_NUM, activation= 'softmax'),
])
model.summary()

BATCH_SIZE = 64
NUM_EPOCHS = 10

cpt_path = 'data/14_text_classifier.hdf5'
checkpoint = tf.keras.callbacks.ModelCheckpoint(cpt_path, monitor='acc', verbose=1, save_best_only= True, mode='max')

model.compile(loss= 'categorical_crossentropy', optimizer='adam', metrics=['acc'])

history= model.fit(partial_x_train, partial_y_train, validation_data= (x_val, y_val),
                   epochs= NUM_EPOCHS, batch_size= BATCH_SIZE, verbose= 1,
                   callbacks=[checkpoint])

epochs = range(1, len(history.history['acc']) + 1)

plt.figure()
plt.plot(epochs, history.history['loss'], 'bo', label='Training loss')
plt.plot(epochs, history.history['val_loss'], 'b', label='Validation loss')
plt.title('Training and validation loss')
plt.xlabel('Epochs')
plt.ylabel('Loss')
plt.legend()
plt.grid()

plt.figure()
plt.plot(epochs, history.history['acc'], 'bo', label='Training acc')
plt.plot(epochs, history.history['val_acc'], 'b', label='Validation acc')
plt.title('Training and validation accuracy')
plt.xlabel('Epochs')
plt.ylabel('Accuracy')
plt.legend()
plt.grid()

results = model.evaluate(test_data, test_first_label)

print('Test loss: {:.4f}'.format(results[0]))
print('Test accuracy: {:.2f} %'.format(results[1]*100))
